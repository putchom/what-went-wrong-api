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
  packId?: string     // "core", "pack.surreal" など
  isActive: boolean
}
```

### 2.4 UserPlan（サブスクプラン）

```ts
UserPlan {
  userId: string
  plan: "free" | "premium"
  expiresAt?: string    // 有効期限（ストア検証と連動）
  updatedAt: string
}
```

---

## 3. REST API

### 3.1 GET /goals

#### 概要
ユーザーの全Goal一覧。

#### レスポンス 200

```json
{
  "goals": [ { ...Goal }, ... ]
}
```

### 3.2 POST /goals

#### 概要

Goal作成。

#### サーバ側ロジック

1. currentGoalCount = SELECT COUNT(*) FROM goals WHERE userId = ?
2. maxGoals を UserPlan から決定
3. currentGoalCount >= maxGoals の場合 → 403 Forbidden

#### リクエストボディ

```json
{
  "title": "読書を1ページ読む",
  "notificationTime": "21:30",
  "notificationEnabled": true
}
```

#### レスポンス 201

```json
{
  "goal": { ...createdGoal }
}
```

### 3.3 PATCH /goals/{goalId}
Goal更新（タイトル、通知等）。
エンタイトルメントに特別な制限なし。

### 3.4 DELETE /goals/{goalId}
Goal削除＋紐づくExcuseEntry削除。
エンタイトルメントに特別な制限なし。

### 3.5 GET /excuse-templates

#### クエリパラメータ（任意）
- `packId`： `core` / `pack.surreal` / `pack.sf` 等
（指定なしの場合、利用可能なすべてのテンプレ）

#### サーバー側ロジック（利用可能なテンプレの絞り込み）

1. UserPlan から canUsePremiumTemplates を判定
2. UserAddOnPurchase から sku一覧を取得
3. 取得条件：
  - Free:
    - `packId = "core"` のテンプレのみ
  - Premium:
    - `packId = "core"` or `packId IS NULL` のテンプレ

#### レスポンス 200

```json
{
  "templates": [
    { "id": "gravity-strong", "text": "今日は重力が強かった", "packId": "core" }
  ]
}
```

### 3.6 GET /goals/{goalId}/excuses

#### 役割

Goalごとの過去ログ取得。

#### クエリ例

- `from` (optional) — `"YYYY-MM-DD"`
- `to` (optional) — `"YYYY-MM-DD"`

#### サーバー側ロジック（Freeの保存期間制限）

- `UserPlan.logRetentionDays`を見て、例：Freeの場合 `logRetentionDays = 30`
- クエリ条件に `date >= today - 30` を強制追加（`from` がそれより後ならそれを優先）

#### レスポンス 200

```json
{
  "excuses": [
    {
      "id": "e1",
      "goalId": "g1",
      "date": "2025-11-28",
      "excuseText": "今日は重力が強かった",
      "templateId": "gravity-strong",
      "createdAt": "2025-11-28T13:00:00Z",
      "updatedAt": "2025-11-28T13:00:00Z"
    }
  ]
}
```

### 3.7 GET /goals/{goalId}/excuses/today

- (goalId, today) の ExcuseEntryを1件返す
- なければ 404 ExcuseEntryNotFoundForToday

（課金制御は特になし）

### 3.8 POST /goals/{goalId}/excuses

#### 今日を含む任意の日付の言い訳保存（upsert）。

#### リクエストボディ

```json
{
  "date": "2025-11-28",
  "excuseText": "今日は重力が強かった",
  "templateId": "gravity-strong"
}
```

#### サーバー側ロジック

- `(userId, goalId, date)` で既存レコードがあれば更新、なければ作成
- `templateId` が指定された場合、それがユーザーに利用可能なテンプレかチェック
  - 利用不可なら 403 Forbidden

#### レスポンス

- 新規時 201 Created
- 更新時 200 OK

```json
{
  "excuse": { ...ExcuseEntry }
}
```

### 3.9 PATCH /excuses/{excuseId}

- `excuseText` / `templateId` を更新
- `templateId` 変更時は利用可能テンプレかチェック

### 3.10 DELETE /excuses/{excuseId}

- ExcuseEntry削除 → その日が「できた扱い」に戻る
- 課金制御なし

### 3.11 GET /me/plan
現在のプランと権限を返す。

```json
{
  "plan": "free",
  "entitlements": {
    "maxGoals": 3,
    "logRetentionDays": 30,
    "canUseAiExcuse": false,
    "canUsePremiumTemplates": false
  }
}
```

### 3.12 POST /me/plan

#### 概要

プラン変更

#### リクエスト

```json
{
  "plan": "premium"
}
```

#### レスポンス 200

```json
{
  "plan": "premium",
  "entitlements": {
    "maxGoals": 100,
    "logRetentionDays": null,
    "canUseAiExcuse": true,
    "canUsePremiumTemplates": true
  }
}
```

### 3.13 POST /ai-excuse

#### 概要

AIで言い訳候補を生成。

#### 前提条件

- `UserPlan.canUseAiExcuse == true`
  - → false の場合 `403 Forbidden`

#### リクエスト

```json
{
  "goalId": "g1",
  "date": "2025-11-28",
  "tone": "surreal",     // "surreal" | "philosophical" | "casual" など任意
  "context": "今日は本を開いたが、SNSを見てしまった"
}
```

#### レスポンス 200

```json
{
  "candidates": [
    "今日はページより通知の方が光って見えた",
    "活字よりもタイムラインが呼んでいた"
  ]
}
```

※ ここでは保存せず、選ばれたものを /excuses にPOSTしてもらう想定。

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

## 6. 課金ルール

| 機能            | free                 | premium            |
| ------------- | -------------------- | ------------------ |
| Goal数         | 最大 3                 | 無制限（実運用上は100など上限可） |
| ログ保存期間（API返却） | 直近30日                | 無制限                |
| テンプレ数（coreのみ） | `packId = "core"` のみ | 全テンプレ（＋購入パック）      |
| AI言い訳生成       | 不可                   | 可能                 |
| 月次レポートAPI     | 単発課金 or プレミアム内包      | プレミアム内包 or 割引      |

## 7. 実装メモ

- どのエンドポイントでも最初に
  - userId をコンテキストに抽出
  - UserPlan を読み込んで entitlements を構築
- entitlements を使って：
  - maxGoals 超過チェック
  - logRetentionDays による ExcuseEntry の絞り込み
  - canUsePremiumTemplates / 購入済パックによるテンプレフィルタ
  - canUseAiExcuse によるAI API制御

これらはミドルウェア or サービスレイヤーで共通化すると保守楽。