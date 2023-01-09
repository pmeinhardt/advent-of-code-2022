extern crate pathfinding;

use std::error::Error;
use std::io::{self, BufRead, Lines};
use std::ops::Index;

use pathfinding::prelude::bfs;

// The approach taken here is a 2d shortest path search with a time dimension.
//
// When computing the possible successors to a position `(x, y)` at time `t`,
// we check which neighbors `(x+1, y), …` are actually valid moves by
// consulting the map state at `t+1`.
//
//     /|/|/|/|/|
//    / / / / / |
// t0 | | | | | | -> t
//    | | | | | |
//    |/|/|/|/|/
//
// If the targeted field is empty at `t+1`, it is a valid move.

#[allow(non_camel_case_types)]
type ts = usize;

#[allow(non_camel_case_types)]
type coord = i8;

type Pos = (coord, coord); // x, y
type Mov = (coord, coord); // x, y

const DIRECTIONS: [Mov; 4] = [
    ( 0, -1), // up
    ( 1,  0), // right
    ( 0,  1), // down
    (-1,  0), // left
];

#[derive(Debug, PartialEq)]
enum Terrain {
    Flat,
    Wall,
}

#[derive(Debug)]
struct Map {
    grid: Vec<Vec<Terrain>>,
    size: (usize, usize),
}

impl Map {
    fn width(&self) -> usize {
        self.size.0
    }

    fn height(&self) -> usize {
        self.size.1
    }
}

impl Index<(usize, usize)> for Map {
    type Output = Terrain;

    fn index(&self, (x, y): (usize, usize)) -> &Self::Output {
        &self.grid[y][x]
    }
}

impl Index<(coord, coord)> for Map {
    type Output = Terrain;

    fn index(&self, (x, y): (coord, coord)) -> &Self::Output {
        &self.grid[u(y)][u(x)]
    }
}

#[derive(Clone, Debug)]
struct Blizzard(Pos, Mov);

#[derive(Debug)]
struct MeteoMap {
    blizzards: Vec<Blizzard>,
}

impl MeteoMap {
    fn step(&self, map: &Map) -> Self {
        let blizzards: Vec<_> = self.blizzards.iter().map(|Blizzard(pos, mov)| {
            let mut nextpos = (pos.0 + mov.0, pos.1 + mov.1);

            if map[nextpos] == Terrain::Wall {
                loop {
                    let moved = (nextpos.0 - mov.0, nextpos.1 - mov.1);
                    if map[moved] == Terrain::Wall { break; }
                    nextpos = moved;
                }
            }

            Blizzard(nextpos, *mov)
        }).collect();

        Self {
            blizzards,
        }
    }

    fn safe(&self, x: coord, y: coord) -> bool {
        !self.blizzards.iter().any(|Blizzard(pos, _)| pos.0 == x && pos.1 == y)
    }
}

fn i(x: usize) -> coord {
    coord::try_from(x).unwrap()
}

fn u(x: coord) -> usize {
    usize::try_from(x).unwrap()
}

fn parse<B: BufRead>(lines: &mut Lines<B>) -> (Map, MeteoMap) {
    let mut blizzards: Vec<Blizzard> = Vec::new();
    let mut grid: Vec<Vec<Terrain>> = Vec::new();
    let mut width: usize = 0;

    for (y, line) in lines.enumerate() {
        let mut row: Vec<Terrain> = Vec::new();

        for (x, ch) in line.unwrap().chars().enumerate() {
            let pos: Pos = (i(x), i(y));

            match ch {
                '#' => row.push(Terrain::Wall),
                _ => row.push(Terrain::Flat),
            }

            match ch {
                '^' => { blizzards.push(Blizzard(pos, DIRECTIONS[0])); },
                '>' => { blizzards.push(Blizzard(pos, DIRECTIONS[1])); },
                'v' => { blizzards.push(Blizzard(pos, DIRECTIONS[2])); },
                '<' => { blizzards.push(Blizzard(pos, DIRECTIONS[3])); },
                _ => {},
            }

        }

        width = width.max(row.len());
        grid.push(row);
    }

    let size: (usize, usize) = (width, grid.len());

    return (Map { grid, size }, MeteoMap { blizzards });
}

fn main() -> Result<(), Box<dyn Error>> {
    let stdin = io::stdin();

    let (map, first) = parse(&mut stdin.lines());

    let mut start: Option<Pos> = None;
    let mut end: Option<Pos> = None;

    let w = map.width();
    let h = map.height();

    for x in 0..w {
        if map[(x, 0)] == Terrain::Flat {
            start = Some((i(x), 0));
        }

        if map[(x, h - 1)] == Terrain::Flat {
            end = Some((i(x), i(h - 1)));
        }
    }

    let mut frames: Vec<MeteoMap> = Vec::new();
    frames.push(first);

    let s = start.expect("no start position");
    let e = end.expect("no end position");

    let iw = i(w);
    let ih = i(h);

    let successors = |&((x, y), t): &(Pos, ts)| {
        // Generate missing (blizzard) map states up to timestamp t
        while frames.len() <= t {
            let next = frames.last().unwrap().step(&map);
            frames.push(next);
        }

        let meteo = &frames[t];

        let mut succ = vec![];
        let tn = t + 1;

        for (dx, dy) in DIRECTIONS {
            let pos = (x + dx, y + dy);

            // Check map bounds
            if pos.0 < 0 || pos.1 < 0 { continue; }
            if pos.0 >= iw || pos.1 >= ih { continue; }

            // Check move is valid
            if map[pos] != Terrain::Flat { continue; }
            if !meteo.safe(x, y) { continue; }

            succ.push((pos, tn)); // move
        }

        if meteo.safe(x, y) {
            succ.push(((x, y), tn)); // stay
        }

        succ
    };

    let success = |&(pos, _): &(Pos, ts)| pos == e;

    let result = bfs(&(s, 0), successors, success);
    let path = result.expect("no path found");

    println!("{:?}", path.len() - 1);

    Ok(())
}
