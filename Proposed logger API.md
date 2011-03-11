LoggerForKeyPath("info.coord.1.proxy.2")
LogAtKeyPath("info.coord.1.proxy.2", "This is an info message")

DeclareOutputPipe("info\b(\..*)", os.Stdout)
DeclareOutputSocket("info\b(\..*)", "127.0.0.1:934234")

ReadConfiguration("log_config")
    {
        "info\b(\..*)": stdout,
        "fatal\b(\..*)": 127.0.0.1:23421,
    }
