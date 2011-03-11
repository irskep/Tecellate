Init(io.Writer)

Logger("info.coord.1.proxy.2")
Log("info.coord.1.proxy.2", "This is an info message")

Redirect("info\b(\..*)", io.Writer)

ReadConfiguration("log_config")
    {
        "info\b(\..*)": stdout,
        "fatal\b(\..*)": 127.0.0.1:23421,
    }
