from proto_compile import compile_go


def rules():
    return [*compile_go.rules()]
