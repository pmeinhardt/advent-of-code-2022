extern crate pathfinding;

use std::io;
use pathfinding::prelude::bfs;

static A: u32 = 'a' as u32;
static Z: u32 = 'z' as u32;

fn h(ch: char) -> u8 {
    return ((ch as u32).clamp(A, Z) - A) as u8;
}

#[derive(Clone, Debug, Eq, Hash, Ord, PartialEq, PartialOrd)]
struct Square(usize, usize);

impl Square {
  fn successors(&self, grid: &Vec<Vec<u8>>) -> Vec<Square> {
    let &Square(x, y) = self;
    let elevation = grid[y][x];

    let max_x = grid[0].len();
    let max_y = grid.len();

    let mut succ: Vec<Square> = vec!();

    if x > 0 && grid[y][x - 1] <= elevation + 1 {
        succ.push(Square(x - 1, y)); // left
    }

    if x < max_x - 1 && grid[y][x + 1] <= elevation + 1 {
        succ.push(Square(x + 1, y)); // right
    }

    if y > 0 && grid[y - 1][x] <= elevation + 1 {
        succ.push(Square(x, y - 1)); // up
    }

    if y < max_y - 1 && grid[y + 1][x] <= elevation + 1 {
        succ.push(Square(x, y + 1)); // down
    }

    return succ;
  }
}

fn main() {
    let mut grid: Vec<Vec<u8>> = Vec::new();

    let mut starts: Vec<Square> = Vec::new();
    let mut end: Option<Square> = None;

    let stdin = io::stdin();

    for (y, line) in stdin.lines().enumerate() {
        let mut row: Vec<u8> = Vec::new();

        for (x, ch) in line.unwrap().chars().enumerate() {
            match ch {
                'S' | 'a' => {
                    starts.push(Square(x, y));
                    row.push(h('a'));
                },
                'E' => {
                    end = Some(Square(x, y));
                    row.push(h('z'));
                },
                _ => {
                    row.push(h(ch));
                },
            }
        }

        grid.push(row);
    }

    let e = end.expect("no end marker");
    let mut lengths: Vec<usize> = Vec::new();

    for s in starts {
        let result = bfs(&s, |p| p.successors(&grid), |p| *p == e);

        match result {
            Some(path) => lengths.push(path.len() - 1),
            None => (),
        }
    }

    println!("{}", lengths.iter().min().expect("no option found"));
}
