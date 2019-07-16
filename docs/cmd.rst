Binary options
==============

This section describes the usage of tranqap command. General usage of the
binary:

.. code:: shell

    tranqap [global flags] [subcommand [subcommand flags]]

Check the following subsections for details.

Global flags
------------

The supported global flags are:

-c string
~~~~~~~~~

Path to config file. (default "config.json")

-l string
~~~~~~~~~

Path to log file. tranqap will not generate additional log file, unless
this option is supplied.

-h
~~

Prints help message

Subcommands
-----------

Subcommands are limited set of feature which are not suitable for the
tranqap shell.

init
----

Generates sample configuration file. Can work with -c flag.

*Example:*

.. code:: shell

    $ tranqap -c config.json init

Creates sample config named config.json in current working directory.
