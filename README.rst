===
foe
===

foe: find, open, edit

``foe`` replaces the two-step process of searching for a file using ``find``
and then opening it with vim.

Running ``foe`` results in finding files and then opening them in a vim session.

All searches are recursive and case-insensitive.

Only text files will be opened; binary files will be ignored.

If more than three files match the query, the matching files will be listed
and vim will not be launched. You may then refine your search terms.

Usage::

    # Find and open any files with "foo" in the name:
    foe foo

    # Find and open any Python files with "foo" in the name:
    foe foo .py 

License
=======

BSD
