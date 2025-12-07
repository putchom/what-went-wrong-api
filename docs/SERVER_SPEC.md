# サーバーサイドAPI 要件定義書

## 1. システム概要

### 1.1 コンセプト
アプリ「なぜできなかったのか」のバックエンドAPI。

- ユーザーは複数の Goal（目標） を持てる
- 各Goalに対し、「できなかった日」だけ ExcuseEntry（言い訳）を保存
- 「できた日」はサーバーに保存しない（存在しない＝できた扱い）
- Goalごとの「今日の状態」と「過去の言い訳一覧」をクライアントが表示するためのAPIを提供

### 1.2 想定クライアント
- iOS / Android
- 1ユーザー = 1アカウント

### 1.3 認証
- 全APIは認証必須
- `Authorization: Bearer <token>`
- トークンから `userId` を復元できる前提

---

## 2. データモデル

### 2.1 Goal

```ts
Goal {
  id: string
  userId: string
  title: string
  notificationTime?: string
  notificationEnabled: boolean
  order: number
  createdAt: string
  updatedAt: string
}
```

### 2.2 ExcuseEntry

```ts
ExcuseEntry {
  id: string
  userId: string
  goalId: string
  date: string
  excuseText: string
  templateId?: string
  createdAt: string
  updatedAt: string
}
```

### 2.3 ExcuseTemplate

```ts
ExcuseTemplate {
  id: string
  text: string
  isActive: boolean
}
```

---

## 3. REST API

### 3.1 GET /goals
Goal一覧取得

### 3.2 POST /goals
Goal作成

### 3.3 PATCH /goals/{goalId}
Goal更新

### 3.4 DELETE /goals/{goalId}
Goal削除（ExcuseEntryも削除）

---

### 3.5 GET /excuse-templates
テンプレ言い訳取得

---

### 3.6 GET /goals/{goalId}/excuses
過去のExcuseEntry取得  
`?from=YYYY-MM-DD&to=YYYY-MM-DD` で期間指定

### 3.7 GET /goals/{goalId}/excuses/today
今日のExcuseEntry取得  
なければ 404

### 3.8 POST /goals/{goalId}/excuses
ExcuseEntry作成（upsert）

### 3.9 PATCH /excuses/{excuseId}
ExcuseEntry更新

### 3.10 DELETE /excuses/{excuseId}
ExcuseEntry削除（その日は成功扱いに戻す）

---

## 4. バリデーション

- Goal.title：1〜200文字  
- ExcuseEntry.excuseText：1〜500文字  
- `(userId, goalId, date)` はユニーク  
- date は `"YYYY-MM-DD"`

---

## 5. 非機能要件

- レスポンス 200ms以内（目安）
- シンプルなDB構成で開始可
- 構造化ログ推奨
