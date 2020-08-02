from .defaults import DefaultScope

class AppScope(DefaultScope):

    @property
    def kind(self):
        return "app"
