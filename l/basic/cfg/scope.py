from typing import Union, Dict
from collections.abc import Iterable


class ConfigScope(object):

    def load(self):
        """Return dictionary with the loaded config values."""



class ConfigScopeSet(set):
    
    known = ("class", "module", "app", "user", "os_cfg", "os_env", "cli")

    def __init__(self, scope_names: Iterable):
        names = set()
        for name in scope_names:
            if name not in self.known:
                raise Exception("Unknown scope: {}".format(name))
            names
        self.values = set(values)


class ConfigValues(object):
    """
    Configuration values can be attached depending on the context and scope.
    """
    def __init__(self, scope: Iterable, defaults=None):
        self.all_scopes = ConfigScopeSet(scope) if not isinstance(scope, ConfigScopeSet) else scope 
        self.defaults = defaults
        self.values = dict()
    
    def load(self):
        """Load configuration values and apply them in the increased order of scope."""
        if scope in ConfigScopeSet.known:
            if scope not in self.all_scopes:
                continue
            # load this scope
            # scope_values = load...
            scope_values = {}
            self.values.update(scope_values)

    def load_in_scope(self, scope_name):



def get_config(name: str, scope: Iterable = None, defaults=None) -> ConfigValues:
    """Get configuration"""
    if scope is None:
        # try to guess the scope by inspecting the context
        pass
    return
