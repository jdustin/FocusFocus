# FocusFocus

A simple tool to help focus on getting stuff done and not procrastinating by looking at sites that waste too much time.

It works by adding or removing a list of sites to your computers hosts file and giving them an address of 127.0.0.1. 

This has the effect of blocking your access to the sites unless you go to the effort of running the focusfocus or manually changing your hosts file. Hopefully this will be annoying enough to just keep working.

- FocusFocus works for both linux and windows (haven't tried a mac). 
- FocusFocus requires sudo/"run as administrator" access to be able to alter the hosts file.
- The hosts file is system wide so obviously this will only be suitable for a system that other people arn't trying to use.


## Getting

Currently FocusFocus has a simple command line util in the focuscmd directory (todo: add gui/web page access)

##Usage

(with sudo or a windows cmd with admin access)

To Focus and waste less time:
Add site domain names to be restricted into a text file, one site per line and run focuscmd with the focus parameter e.g.

focuscmd focus sitestolimit.txt

To Relax, remove any sites added by FocusFocus, run focuscmd with the relax param. e.g.

focuscmd relax
