from .defaults import DefaultScope

class ModuleScope(DefaultScope):

    @property
    def kind(self):
        return "mod"

    def load(self):
        """Load values from other sources."""
        