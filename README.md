# Vertical Scaling Function

OCI インスタンスのシェイプを指定した二つのシェイプ間で切り替えます。 \
以降二つのシェイプをアクティブ・インアクティブなシェイプと呼称します。

# 環境変数

以下の環境変数をファンクションの構成に設定してください。

```
INSTANCE_ID: インスタンスの OCID
ACTIVE_OCPU: アクティブ時の OCPU 数
ACTIVE_MEMORY: アクティブ時のメモリ量 (GiB)
INACTIVE_OCPU: インアクティブ時の OCPU 数
INACTIVE_MEMORY: インアクティブ時のメモリ量 (GiB)
```

# 導入方法

記事作成中