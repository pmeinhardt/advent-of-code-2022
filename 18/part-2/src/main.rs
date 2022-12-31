extern crate pathfinding;

use std::cmp::{max,min};
use std::collections::HashSet;
use std::io;

use pathfinding::prelude::astar;

type Point = (i8, i8, i8); // x, y, z

const DIRECTIONS: [(i8, i8, i8); 6] = [
    ( 1,  0,  0),
    (-1,  0,  0),
    ( 0,  1,  0),
    ( 0, -1,  0),
    ( 0,  0,  1),
    ( 0,  0, -1),
];

fn neighbors(p: Point) -> [Point; 6] {
    return DIRECTIONS.map(|v| (p.0 + v.0, p.1 + v.1, p.2 + v.2));
}

fn distance(a: Point, b: Point) -> u8 {
    return a.0.abs_diff(b.0) + a.1.abs_diff(b.1) + a.2.abs_diff(b.2);
}

fn main() {
    let mut cubes: HashSet<Point> = HashSet::new();

    let stdin = io::stdin();

    for line in stdin.lines() {
        let coords: Vec<i8> = line.unwrap()
                                  .split(',')
                                  .map(|s| s.parse::<i8>().unwrap())
                                  .collect();

        match coords[..] {
            [x, y, z] => { cubes.insert((x, y, z)); },
            _ => panic!("Invalid input")
        }
    }

    let mut faces = cubes.len() * 6;

    let sample = cubes.iter().next().expect("no cubes");

    let mut pmin = *sample;
    let mut pmax = *sample;

    for c in &cubes {
        for n in neighbors(*c) {
            if cubes.contains(&n) {
                faces -= 1;
            }
        }

        pmin = (min(pmin.0, c.0), min(pmin.1, c.1), min(pmin.2, c.2));
        pmax = (max(pmax.0, c.0), max(pmax.1, c.1), max(pmax.2, c.2));
    }

    let outside = (pmax.0 + 1, pmax.1 + 1, pmax.2 + 1);

    for x in pmin.0..=pmax.0 {
        for y in pmin.1..=pmax.1 {
            for z in pmin.2..=pmax.2 {
                let p = (x, y, z);

                if cubes.contains(&p) {
                    continue;
                }

                let result = astar(&p,
                                   |&p| neighbors(p).into_iter().filter(|n| !cubes.contains(n)).map(|p| (p, 1)),
                                   |&p| distance(p, outside),
                                   |&p| p == outside);

                if result.is_none() {
                    for n in neighbors(p) {
                        if cubes.contains(&n) {
                            faces -= 1;
                        }
                    }
                }
            }
        }
    }

    println!("{}", faces)
}
