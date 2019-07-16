Installation
============

The application is written on Go, which makes it relatively portable. It
is developed and tested for Linux, but it is supposed to work on other
OSes too, thanks to the nature of GoLang.

Binary installation
-------------------

For each release rpm and deb packages are provided. Head to
`releases <https://github.com/tdimitrov/tranqap/releases>`__ tab of the
project and download the package for your distribution.

Currently there are two options:

-  RPM for Fedora.
-  DEB for Ubuntu.

The packages are tested on the latest stable releases of the
distributions. For Fedora this means the current stable release, for
Ubuntu - latest standard and latest LTS releases.

Installation from source
------------------------

tranqap is written in Go, so the Go distribution should be installed.
Instructions for the installation can be found
`here <https://golang.org/doc/install>`__.

After that:

.. code:: bash

    go get  github.com/tdimitrov/tranqap
    go install  github.com/tdimitrov/tranqap

At this point tranqap should be installer in ``$GOPATH/bin``. This path
should be added to system path.
