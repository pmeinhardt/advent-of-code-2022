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

    for c in &cubes {
        let neighbors = DIRECTIONS.map(|v| (c.0 + v.0, c.1 + v.1, c.2 + v.2));

        for n in neighbors {
            if cubes.contains(&n) {
                faces -= 1;
            }
        }
    }

    println!("{}", faces)
}
