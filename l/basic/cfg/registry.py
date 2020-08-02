"""Key-value configuration registry with multiple scopes.

Context is a list of multiple scope instances.

Lookup is a 2-step process: (key, context) -> scopes (sources) -> value

Registry allows idempotent collection of runtime configuration
values in multiple contexts from many optional sources. Values are
never overwritten.

It is possible to dump registry runtime content for debugging purposes.

If needed keys can be prefixed to achieve namespase-like logical separation.
"""


class ScopeOrder(set):

    known = (
        "cls",     # Class-level
        "mod",     # Module-level
        "app",     # Application-level
        "usr",     # User-level config file
        "os",      # OS-wide config file
        "env",     # OS ENV vars
        "cli",     # command-line args
    )

    @classmethod
    def is_valid(cls, scope):
        return scope in cls.known

    @classmethod
    def get_upstream(cls, start_scope):
        """Return list of upstream scopes starting with scope."""
        if not cls.is_valid(start_scope):
            raise Exception("Unknown scope: {}".format(start_scope))
        r = list(cls.known)
        return r[r.index(start_scope):]



class HierarchicalRegistry:
    """Runtime hierarchical Keyvalue config registry."""

    def __init__(self):
        self.contexts = dict()

    def add(self, context_name: str, scopes: list):
        self.contexts[context_name] = scopes

    def get(self, context_name, key):
        scopes = self.contexts[context_name]
        for scope in scopes:
            if key in scope.values:
                # scheme validation is taken place at load-time.
                return scope.values.get(key)
        raise KeyError(f"Failed to find key '{key}'in any sources.")

    def dump(self):
        """Dump registry runtime content for debugging purposes"""


config_registry = KeyvalueRegistry()


def get_config(context, sources=None, defaults=None):
    """
    Get config base on the context key.

    context:

        "mod:l.basic.cli.parser" corresponds to module scope config anchor.

        "cls:l.basic.cli.parser.Parser" corresponds to class scope config anchor.

    sources:

        "yaml" - at module-level scope will look for "<modfile>.yaml" file sibling
        near "<modfile>.py". This will fail if not found.

        "toml" - same as with yaml.

    """
    scope_names =
