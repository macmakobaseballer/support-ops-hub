plugin "aws" {
  enabled = true
  version = "0.40.0"
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

config {
  format = "compact"
}

# M0〜M7 は modules/* がスケルトン（resource 未実装）のため、宣言済み変数が
# 未使用と判定される。M8 で実リソースを実装した時点でこのルールを再有効化する。
rule "terraform_unused_declarations" {
  enabled = false
}

# 同様に、骨格モジュールに terraform { required_version } を強制しない。
# envs/{dev,stg,prod}/main.tf で provider 制約を一元管理しているので十分。
# M8 で各モジュールが実 resource を持つ段階で再有効化する。
rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}
