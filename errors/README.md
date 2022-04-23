# Errors Package

Lib designed to facilitate and standardize how error handling should be done inside Trivela's systems.

It produces `CustomError`, that encodes useful information about a given error.

It's supposed to flow within the application in detriment of the the default golang error since its `Kind` and `Code` attributes are the keys 
to express its semantic and uniqueness, respectively.

It should be generated once by the peace of code that found the error (because it's where we have more context about the error),
and be by passed to the upper layers of the application.
