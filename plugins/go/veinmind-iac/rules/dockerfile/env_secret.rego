package brightMirror.dockerfile

import future.keywords.in
import data.common

risks[res] {
    inner := input[_]
    inner.Cmd == "env"
    some val in inner.Value
        val in ["passwd","password","pass","secret","key","access","api_key","apikey","token","tkn"]
    res := common.result(inner, "DF-012")
}