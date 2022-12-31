use std::cmp::{max,min};
use std::collections::HashSet;
use std::io;

type Point = (i8, i8, i8); // x, y, z

const DIRECTIONS: [(i8, i8, i8); 6] = [
    ( 1,  0,  0),
    (-1,  0,  0),
    ( 0,  1,  0),
    ( 0, -1,  0),
    ( 0,  0,  1),
    ( 0,  0, -1),
];

fn neighbors(p: &Point) -> [Point; 6] {
    return DIRECTIONS.map(|v| (p.0 + v.0, p.1 + v.1, p.2 + v.2));
}

fn within(p: &Point, pmin: &Point, pmax: &Point) -> bool {
    return p.0 >= pmin.0 && p.1 >= pmin.1 && p.2 >= pmin.2
        && p.0 <= pmax.0 && p.1 <= pmax.1 && p.2 <= pmax.2;
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

    // Compute bounding box

    let sample = cubes.iter().next().expect("no cubes");

    let mut pmin = *sample;
    let mut pmax = *sample;

    for c in &cubes {
        pmin = (min(pmin.0, c.0), min(pmin.1, c.1), min(pmin.2, c.2));
        pmax = (max(pmax.0, c.0), max(pmax.1, c.1), max(pmax.2, c.2));
    }

    // Expand bounding box by 1 unit in each direction

    pmin = (pmin.0 - 1, pmin.1 - 1, pmin.2 - 1);
    pmax = (pmax.0 + 1, pmax.1 + 1, pmax.2 + 1);

    // Flood area around the body, within bounding box limits

    let mut closed: HashSet<Point> = HashSet::new();
    let mut open: HashSet<Point> = HashSet::from([pmin, pmax]);

    while open.len() > 0 {
        let mut next: HashSet<Point> = HashSet::new();

        for p in open {
            for n in neighbors(&p) {
                if cubes.contains(&n) { continue; }
                if closed.contains(&n) { continue; }
                if !within(&n, &pmin, &pmax) { continue; }
                next.insert(n);
            }

            closed.insert(p);
        }

        open = next;
    }

    // Count surfaces where flooded cubes touch cubes of the scanned body

    let mut faces = 0;

    for c in closed {
        for x in neighbors(&c) {
            if cubes.contains(&x) {
                faces += 1;
            }
        }
    }

    println!("{}", faces)
}
