# Particles

"Particles" is a Go implementation of the algorithm presented in [this video](https://youtu.be/0Kx4Y9TVMGg). Basically, the user can create a set of populations defined by a color, each population has a quantity of individuals. Each individual has a 2d position in a constrained space. The user can then create rules affecting all individuals of a population $\alpha$. A rule specifies that individuals of a population $\beta$ affects individuals of this population $\alpha$ (it can be the same) with a force *f*, and with an effect range. If *f* > 0 then individuals of $\beta$ repel individuals of $\alpha$, and attract them otherwise. It appears that even a few rules can create interesting results, and complex behaviours can emerge from this algorithm even though not strictly defined in the code. I found this idea inspiring, and this is my attempt to implement it !

## Live demo

The result of this code can be viewed here: [salicorne.github.io/particles](https://salicorne.github.io/particles/)

## Roadmap

- [x] Core algorithm implementation
- [x] User interface with bootstrap
- [ ] Definition of rules from the UI
- [ ] Random seed definition
- [ ] Export and import of settings
- [ ] Fullscreen mode