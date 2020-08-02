class DefaultScope:
    """
    Default constant scope without any anchoring location.
    """

    def __init__(self, name, defaults=None):
        self.name = name
        self.values = dict()
        if defaults:
            self.values.update(defaults)
            # Defaults can be modified/extended before application run.

    @property
    def kind(self):
    	return "defaults"

    def load(self):
        """Load config keyvalue entries."""
        # Defaults are defined at construction time

    def __str__(self):
        return f"{self.kind}:{self.name}"
