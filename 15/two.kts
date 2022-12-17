/**
 * Approach:
 *
 * We have an area, 4000000x4000000, and are looking for a point that is not
 * covered by any of our 32 sensors. We know there is exactly one spot.
 *
 * Because of the "diamond" shape of the sensor coverage areas, such a single
 * uncovered spot can only appear either in one of the four corners of our
 * search area or on an intersection between two outside edges of two sensor
 * coverage areas.
 *
 *   +--- … -----
 *   |.## … ##.##
 *   |### … #####
 *
 */

data class Point(val x: Int, val y: Int) {
    fun offset(other: Point): Pair<Int, Int> {
        return Pair(other.x - this.x, other.y - this.y)
    }

    fun distance(other: Point): Int {
        val (dx, dy) = this.offset(other)
        return kotlin.math.abs(dx) + kotlin.math.abs(dy)
    }

    override fun toString(): String = "($x,$y)"
}

class LineSegment(a: Point, b: Point) {
    val start: Point
    val end: Point

    val slope: Int // we only have slopes of 1 and -1
    val line: Line

    init {
        if (a.x < b.x) {
            start = a
            end = b
        } else if (b.x < a.x) {
            start = b
            end = a
        } else {
            throw RuntimeException("Invalid points for line segment: $a and $b")
        }

        slope = (end.y - start.y) / (end.x - start.x)
        val y0 = start.y - (start.x * slope)

        line = Line(y0, slope)
    }

    fun contains(point: Point): Boolean {
        val (px, py) = point
        val (sx, sy) = start
        val (ex, ey) = end

        if (px < sx) return false
        if (px > ex) return false

        val dx = px - sx

        return py - slope * dx == sy
    }

    fun intersect(other: LineSegment): Point? {
        val point = line.intersect(other.line)
        if (point == null) return null

        if (!this.contains(point)) return null
        return point
    }

    override fun toString(): String = "$start <-> $end"
}

class Line(val y0: Int, val slope: Int) {
    fun intersect(other: Line): Point? {
        if (other.slope == slope) return null // ignore parallel lines

        val x = (y0 - other.y0) / (other.slope - slope)
        val y = y0 + slope * x

        return Point(x, y)
    }

    override fun toString(): String = "y = $y0 + $slope * x"
}

class Sensor(val position: Point, val beacon: Point) {
    val radius: Int = this.beacon.distance(this.position)
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

fun outsideEdgesOf(sensor: Sensor): Iterable<LineSegment> {
    val (x, y) = sensor.position
    val r = sensor.radius

    val top = Point(x, y - (r + 1))
    val bottom = Point(x, y + (r + 1))
    val left = Point(x - (r + 1), y)
    val right = Point(x + (r + 1), y)

    return listOf(
        LineSegment(top, right),
        LineSegment(bottom, right),
        LineSegment(left, top),
        LineSegment(left, bottom)
    )
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


// Main

require(args.size > 0)

val limit = args.first().toInt()
val (sensors, beacons) = readInput()

val corners = listOf(Point(0, 0), Point(0, limit), Point(limit, 0), Point(limit, limit))

val points = mutableListOf<Point>()

points.addAll(corners)

val edges = sensors.flatMap { outsideEdgesOf(it) }
val intersections = mutableSetOf<Point>()

for ((i, j) in pairs(0..(edges.count() - 1))) {
    val a = edges[i]
    val b = edges[j]

    val p = a.intersect(b)

    if (p != null && p.x >= 0 && p.y >= 0 && p.x <= limit && p.y <= limit) {
        intersections.add(p)
    }
}

points.addAll(intersections)

for (point in points) {
    if (!sensors.any { sensor -> point.distance(sensor.position) <= sensor.radius }) {
        println(point.x.toLong() * limit + point.y)
        break
    }
}