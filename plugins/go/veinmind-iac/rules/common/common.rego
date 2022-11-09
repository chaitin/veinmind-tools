package common

result(obj, meta) = result{
    result := {
        "rule": meta_data[meta],
        "risk": {
            "startline": object.get(obj, "startline", object.get(obj, "StartLine", 0)),
            "endline": object.get(obj, "endline", object.get(obj, "EndLine", 0)),
            "filePath": object.get(obj, "filepath", object.get(obj, "Path" , "")),
            "original": object.get(obj, "original", object.get(obj, "Original", "")),
        }
    }
}