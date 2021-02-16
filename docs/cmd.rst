Binary options
==============

This section describes the usage of the tranqap binary. There are three 
type of parameters which can be passed to it - global flas, subcommand
and subcommand flags. They all are optional:

.. code:: shell

    tranqap [global flags] [subcommand [subcommand flags]]

Check the following subsections for details.

Global flags
------------

These flags represent global configuration options for tranqap. The 
supported global flags are:

-c string
~~~~~~~~~

Path to configuration file. (default "config.yaml")

-l string
~~~~~~~~~

Path to log file. tranqap will not generate log file, unless this 
option is supplied. The log file is useful mainly for debugging by 
tranqap developers.

-h
~~

Prints help message

Subcommands
-----------

Subcommands are limited set of feature which are not suitable for the
tranqap shell. At the moment only one subcommand is supported.

init
~~~~

Generates sample configuration file. Can work with -c flag.

*Example:*

.. code:: shell

    $ tranqap -c config.yaml init

Creates sample config named config.yaml in current working directory.
