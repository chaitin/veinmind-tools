class register:
    plugin_dict = {}
    plugin_name = []

    @classmethod
    def register(cls, plugin_name):
        def wrapper(plugin):
            cls.plugin_dict[plugin_name] = plugin
            return plugin
        return wrapper