
# cfg

## scope and location

Scope depends on the context and the "name space" where configuration is applied or relevant to.

Configuration scope does somewhat predetermines location - a concrete place where and how the exact 
configuration values are specified.

The specificity can increase as the scope becomes wider.

    - class specific 
    - module specific
    - application wide
    - user specific
    - OS specific
    - runtime specific:
     - sourced from the DB or remote registry
     - OS envs
     - cli args

## supported formats

Configuration can be declared or stored in different representation formats. Practicality of each depends on the scope.

For example, at module or class level declaration can be made more plain - directly in python structures and objects 
like nested dictionary. When configuration is loaded from a file it can be either YAML or TOML.

## access configuration values

The end result must be convertible into plain python data types or objects.

Often a preferable contract to work with configuration is python object (that can wrap around nested dictionary and walk the values).
