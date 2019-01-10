# 2D vision

Implementing a some game graphic concepts in Go, using [ebiten](https://hajimehoshi.github.io/ebiten/) 2d graphics library. Thanks to [@hajimehoshi](https://twitter.com/hajimehoshi) for support. Examples are hosted at jsgo.io. Latest version here: [https://jsgo.io/kyeett/2d-vision](https://jsgo.io/kyeett/2d-vision)

Heavily inspired by Red Blob Games' article on [2D visibility](https://www.redblobgames.com/articles/visibility/). Also found this excellent wikipedia article on [Visibility polygon](https://en.wikipedia.org/wiki/Visibility_polygon#Angular_sweep))

### Disclaimer

The code in this library is quite messy. Use with care :-)

## Basic ray casting

Simple algorithm. Send rays in all directions. If it intersects with one of the lines (walls or boxes), select the point that is closest to the player and draw a line to there. These lines used as vertices for triangles that are removed from the shadows.

[Demo here](https://pkg.jsgo.io/github.com/kyeett/2d-vision.c62ef4d28cc90c6ee0aa8239ca38031bfa170bd7.js)

![Basic Ray Casting](/doc/basic_ray_casting.gif)

## Smarter ray casting

Calculate the angles to the corners of all objects, and do the ray cast only in that angle.

**Note**: At first, I had a lot of problems with the rays stopping ON the corners. I solved this by sending _two_ lines per corner, with a small offset to the angle in both directions. It works well, but doubles the number of lines. I'm not sure if this is the right approach.

[Demo here](https://pkg.jsgo.io/github.com/kyeett/2d-vision.c62ef4d28cc90c6ee0aa8239ca38031bfa170bd7.js)

![Smart Ray Casting](/doc/smart_ray_casting.gif)

### Resources:

- [Floor](https://opengameart.org/content/even-grey-stone-tile-floor-256px) by [Tiziana](http://www.unbruco.it/offcircle/index_en.html)
