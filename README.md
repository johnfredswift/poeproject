# poeproject
This project was created as a my final year project at University.
CLI tool to analysis dump data from poe.ninja including graphs, most explaination and evaluation of the project is within the project proposal and report as to not infringe on word count and/or plagerism.

This project used data-dump files from poe.ninja/data-dumps on a popular ARPG, Path of Exile, to create anaytlics and analysis for players evaluating the Path of Exile economies.

## Installation

After pulling the files the data-dumps are added to raw-data section, from there the tool can be setup and used using the following sub-commands

'-update': Expects no arguements, This creates the requires file directories and populates them.

'-cont': Expects one arguement, two flags are avalible to choose if the data should be on Hardcore Leagues 'hc' or Temporary Leagues 'st'.
The graph from this is created and put within the graphs folder.

'sing': Expects two types of arguement, firsly a single item out of the entire history of the game, secondly any number of league names of which you want the item price compared on. This has no flags avaliable to also for comparison over standard, temporary, hardcore and softcore leagues alike.

'json':Expects no arguements, Outputs a JSON file with information on all the leagues, this serves as a fast way of accessing the basic information without additional looking outside of the tool.

## Infomation

There is also test files included to help diagnose and ensure proper installition. These are all implimented as GoLang tests.
