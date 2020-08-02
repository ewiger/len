from .defaults import DefaultScope

class ClassScope(DefaultScope):

    @property
    def kind(self):
        return "cls"
