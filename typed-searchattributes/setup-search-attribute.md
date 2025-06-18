# Search Attributes の設定方法

## 登録済みのSearch Attributesを確認する

```sh
temporal operator search-attribute list
```

## 任意のSearch Attributesを登録する

```sh
temporal operator search-attribute create \
    --name YourAttributeName \
    --type Keyword
```

### サンプルで使用しているSearch Attributesを登録する

```sh
temporal operator search-attribute create --name CustomIntField --type Int
```
```sh
temporal operator search-attribute create --name CustomStringField --type Text
```
```sh
temporal operator search-attribute create --name CustomKeywordField --type Keyword
```
```sh
temporal operator search-attribute create --name CustomBoolField --type Bool
```
```sh
temporal operator search-attribute create --name CustomDatetimeField --type Datetime
```
```sh
temporal operator search-attribute create --name CustomDoubleField --type Double
```
