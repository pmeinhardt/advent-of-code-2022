data class Point(val x: Int, val y: Int) {
    fun offset(other: Point): Pair<Int, Int> {
        return Pair(other.x - this.x, other.y - this.y)
    }

    fun distance(other: Point): Int {
        val (dx, dy) = this.offset(other)
        return kotlin.math.abs(dx) + kotlin.math.abs(dy)
    }
}

class Sensor(val position: Point, val beacon: Point) {
    fun radius(): Int = this.beacon.distance(this.position)
}

fun readInput(): Pair<Collection<Sensor>, Collection<Point>> {
    val sensors = mutableListOf<Sensor>()
    val beacons = mutableSetOf<Point>() // use set as multiple sensors may acquire the same beacon

    val regex = Regex("x=(-?[0-9]+), y=(-?[0-9]+):.*x=(-?[0-9]+), y=(-?[0-9]+)")

    while (true) {
        val line = readLine() ?: break
        val match = regex.find(line)!!
        val (x, y, bx, by) = match.destructured

        val position = Point(x.toInt(), y.toInt())
        val beacon = Point(bx.toInt(), by.toInt())

        sensors.add(Sensor(position, beacon))
        beacons.add(beacon)
    }

    return Pair(sensors, beacons)
}

fun coverage(sensor: Sensor, y: Int): IntRange {
    val (sx, sy) = sensor.position
    val offset = kotlin.math.abs(y - sy)
    val radius = sensor.radius()

    if (offset > radius) {
        return IntRange.EMPTY
    }

    val rem = radius - offset
    return IntRange(sx - rem, sx + rem)
}

fun touching(a: IntRange, b: IntRange): Boolean {
    val af = a.first()
    val al = a.last()
    val bf = b.first()
    val bl = b.last()
    return a.contains(bf) || a.contains(bf - 1)
        || a.contains(bl) || a.contains(bl + 1)
        || b.contains(af) || b.contains(af - 1)
        || b.contains(al) || b.contains(al + 1)
}

fun pairs(range: IntRange): Sequence<Pair<Int, Int>> {
    val min = range.first()
    val max = range.last()

    var i = min
    var j = min

    return generateSequence {
        j = j + 1

        if (j > max) {
            i = i + 1
            j = i + 1
        }

        if (i > max - 1) null else Pair(i, j)
    }
}

fun merge(ranges: Collection<IntRange>): Collection<IntRange> {
    val rs = ranges.filter { range -> !range.isEmpty() }.toMutableList()

    while (rs.size > 1) {
        var changed = false

        for ((i, j) in pairs(0..(rs.size - 1))) {
            val a = rs[i]
            val b = rs[j]

            if (touching(a, b)) {
                rs.removeAt(j) // remove higher-index element first!
                rs.removeAt(i)

                val start = kotlin.math.min(a.first(), b.first())
                val end = kotlin.math.max(a.last(), b.last())
                rs.add(IntRange(start, end))

                changed = true
                break
            }
        }

        if (!changed) break // no ranges were merged
    }

    return rs
}

// Main

require(args.size > 0)

val row = args.first().toInt()
val (sensors, beacons) = readInput()

val ranges = merge(sensors.map { sensor -> coverage(sensor, row) })

val covered = ranges.fold(0) { total, range -> total + range.count() }
val bcount = beacons.count { beacon -> beacon.y == row }

println(covered - bcount)
