# temporal-sample
temporalの検証用リポジトリ

-[Temporal Workflow | Temporal Platform Documentation](https://docs.temporal.io/workflows)のサンプルを動かしてみる

## サーバー起動方法

```bash
temporal server start-dev
```

### DBの起動やポートの指定して起動する

```bash
temporal server start-dev --db-filename temporal.db --ui-port 8080
```
