# Silvercord 
Ask Silvercord anything about CUNY on your Discord server! Including: the academic schedule, upcoming or current professors, course material & how to best study for it, or upcoming courses to be opened or closed.

## Demo
Coming Soon

## Purpose
Many students including new-comers tend to have many important questions about the school's academics. Despite the information being searchable online, it would be easier to have everything about each school stored
in one place not only for information, but to be interactive and capable of having a discussion about such a topic. 

# Setting Up
Initialize your NPM environment: </br>
```bash
npm install
``` </br>
Make sure before running, you have received the .env file and placed it into your Silvercord directory. </br>

Nodemon is a separate package that mnust be installed globally.

For Windows the following command should work:
```bash
npm install -g nodemon
``` 
For MacOS/Linux the following command should work: 
```bash
sudo npm install -g nodemon
``` 
Run the bot: </br>
```bash
nodemon
```</br>

### Setting Up (adv)
While I put advance in the title, it's not really that advanced. I have created a flake.nix file that contains all the packages, and even exports the token for you. You can download the Nix package manager and run the following command to install all the dependencies: 
```bash
nix develop
```
This will install all deppendencies (python and nodejs), leaving you in a preconfigured shell that is guaranteed to work with the project on any machine. Check out [this](https://nix.dev/install-nix) link for more information on how to get started with Nix.

## Use Cases

- Silvercord will answer whenever users ask about the academic schedule
- Deliver sentiment analysis on professors at CUNY to help students make a decision on their course schedules
- Inform users of course material to be extra-aware of the courses they will be enrolling into
- Help students prepare the best way to study for a certain course by generating outlines, resources, etc.
- The ability to check for updates on the Coursicle platform, specifically course section seats

## Limitations
