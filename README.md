# Despensa App


## Project Structure

- public - Assets that don't need server side processing.
 - css - 3Party or final css builds.
 - favicon - Favicon and app icon assets.
 - fonts - Custom fonts.
 - images - Static images like the logo or background.
 - js - 3Party or final js builds.

- assets - Assets that need server side processing or private.
 - style - SCSS or CSS app specific that require processing.
 - client - Javascript or Typescript app specific
 - images - Images that will be processed before being accessible
 - favicon - 


## Make Targets

- help: Print the help message (default target)
- watch: Starts [Air](https://github.com/cosmtrek/air) in Live Reload and Tailwindcss and Templ in watch mode and hot-reload
- build: Builds the application in production mode including Tailwindcss and Templ
- run: Runs the production build

