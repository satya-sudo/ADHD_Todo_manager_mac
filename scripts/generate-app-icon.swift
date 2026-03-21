import AppKit
import Foundation

guard CommandLine.arguments.count >= 2 else {
    fputs("usage: generate-app-icon.swift <output-png>\n", stderr)
    exit(1)
}

let outputPath = NSString(string: CommandLine.arguments[1]).expandingTildeInPath
let size = NSSize(width: 1024, height: 1024)
let image = NSImage(size: size)

func color(_ r: CGFloat, _ g: CGFloat, _ b: CGFloat, _ a: CGFloat = 1.0) -> NSColor {
    NSColor(calibratedRed: r / 255.0, green: g / 255.0, blue: b / 255.0, alpha: a)
}

func roundedPath(_ rect: NSRect, radius: CGFloat) -> NSBezierPath {
    NSBezierPath(roundedRect: rect, xRadius: radius, yRadius: radius)
}

image.lockFocus()

NSColor.clear.setFill()
NSRect(origin: .zero, size: size).fill()

let cardRect = NSRect(x: 118, y: 118, width: 788, height: 788)
let cardPath = roundedPath(cardRect, radius: 185)
color(31, 31, 36).setFill()
cardPath.fill()

let cardGlow = NSShadow()
cardGlow.shadowBlurRadius = 42
cardGlow.shadowOffset = NSSize(width: 0, height: -12)
cardGlow.shadowColor = color(0, 0, 0, 0.28)
NSGraphicsContext.saveGraphicsState()
cardGlow.set()
color(33, 33, 39).setFill()
cardPath.fill()
NSGraphicsContext.restoreGraphicsState()

let boltShadow = NSShadow()
boltShadow.shadowBlurRadius = 36
boltShadow.shadowOffset = NSSize(width: 0, height: 0)
boltShadow.shadowColor = color(255, 193, 30, 0.55)

let bolt = NSBezierPath()
bolt.move(to: NSPoint(x: 620, y: 820))
bolt.line(to: NSPoint(x: 404, y: 504))
bolt.line(to: NSPoint(x: 540, y: 504))
bolt.line(to: NSPoint(x: 454, y: 300))
bolt.line(to: NSPoint(x: 676, y: 572))
bolt.line(to: NSPoint(x: 554, y: 572))
bolt.close()

NSGraphicsContext.saveGraphicsState()
boltShadow.set()
let boltGradient = NSGradient(colors: [
    color(255, 245, 160),
    color(255, 211, 62),
    color(255, 180, 23),
])!
boltGradient.draw(in: bolt, angle: -65)
NSGraphicsContext.restoreGraphicsState()

let boltHighlight = NSBezierPath()
boltHighlight.move(to: NSPoint(x: 606, y: 774))
boltHighlight.line(to: NSPoint(x: 450, y: 528))
boltHighlight.line(to: NSPoint(x: 520, y: 528))
boltHighlight.line(to: NSPoint(x: 470, y: 386))
boltHighlight.close()
color(255, 252, 225, 0.82).setFill()
boltHighlight.fill()

let baseShadow = NSShadow()
baseShadow.shadowBlurRadius = 30
baseShadow.shadowOffset = NSSize(width: 0, height: -10)
baseShadow.shadowColor = color(135, 255, 64, 0.20)

let stem = NSBezierPath()
stem.move(to: NSPoint(x: 509, y: 297))
stem.curve(to: NSPoint(x: 509, y: 368), controlPoint1: NSPoint(x: 504, y: 322), controlPoint2: NSPoint(x: 503, y: 348))
stem.curve(to: NSPoint(x: 523, y: 398), controlPoint1: NSPoint(x: 513, y: 381), controlPoint2: NSPoint(x: 519, y: 392))
stem.lineWidth = 14
stem.lineCapStyle = .round
stem.lineJoinStyle = .round
NSGraphicsContext.saveGraphicsState()
baseShadow.set()
color(118, 230, 44).setStroke()
stem.stroke()
NSGraphicsContext.restoreGraphicsState()

func drawLeaf(points: [NSPoint], fillColors: [NSColor], strokeColor: NSColor) {
    let leaf = NSBezierPath()
    leaf.move(to: points[0])
    leaf.curve(to: points[1], controlPoint1: points[2], controlPoint2: points[3])
    leaf.curve(to: points[4], controlPoint1: points[5], controlPoint2: points[6])
    leaf.curve(to: points[0], controlPoint1: points[7], controlPoint2: points[8])
    leaf.close()

    let gradient = NSGradient(colors: fillColors)!
    NSGraphicsContext.saveGraphicsState()
    baseShadow.set()
    gradient.draw(in: leaf, angle: 22)
    NSGraphicsContext.restoreGraphicsState()

    leaf.lineWidth = 3
    strokeColor.setStroke()
    leaf.stroke()
}

drawLeaf(
    points: [
        NSPoint(x: 507, y: 396),
        NSPoint(x: 360, y: 442),
        NSPoint(x: 454, y: 395),
        NSPoint(x: 390, y: 404),
        NSPoint(x: 367, y: 383),
        NSPoint(x: 336, y: 433),
        NSPoint(x: 338, y: 395),
        NSPoint(x: 400, y: 362),
        NSPoint(x: 458, y: 370),
    ],
    fillColors: [color(220, 246, 114), color(148, 225, 66), color(100, 190, 48)],
    strokeColor: color(98, 160, 46, 0.72),
)

drawLeaf(
    points: [
        NSPoint(x: 521, y: 395),
        NSPoint(x: 684, y: 462),
        NSPoint(x: 570, y: 399),
        NSPoint(x: 632, y: 418),
        NSPoint(x: 654, y: 378),
        NSPoint(x: 706, y: 451),
        NSPoint(x: 705, y: 402),
        NSPoint(x: 628, y: 357),
        NSPoint(x: 562, y: 366),
    ],
    fillColors: [color(196, 242, 101), color(128, 214, 53), color(84, 183, 40)],
    strokeColor: color(88, 154, 41, 0.72),
)

func drawVein(from: NSPoint, control1: NSPoint, control2: NSPoint, to: NSPoint) {
    let vein = NSBezierPath()
    vein.move(to: from)
    vein.curve(to: to, controlPoint1: control1, controlPoint2: control2)
    vein.lineWidth = 5
    vein.lineCapStyle = .round
    color(240, 255, 190, 0.65).setStroke()
    vein.stroke()
}

drawVein(
    from: NSPoint(x: 494, y: 401),
    control1: NSPoint(x: 461, y: 415),
    control2: NSPoint(x: 420, y: 424),
    to: NSPoint(x: 374, y: 417),
)

drawVein(
    from: NSPoint(x: 536, y: 403),
    control1: NSPoint(x: 570, y: 423),
    control2: NSPoint(x: 617, y: 436),
    to: NSPoint(x: 672, y: 430),
)

let soil = NSBezierPath()
soil.move(to: NSPoint(x: 380, y: 262))
soil.curve(to: NSPoint(x: 650, y: 262), controlPoint1: NSPoint(x: 442, y: 298), controlPoint2: NSPoint(x: 584, y: 298))
soil.curve(to: NSPoint(x: 608, y: 292), controlPoint1: NSPoint(x: 632, y: 277), controlPoint2: NSPoint(x: 621, y: 290))
soil.curve(to: NSPoint(x: 522, y: 295), controlPoint1: NSPoint(x: 584, y: 304), controlPoint2: NSPoint(x: 546, y: 305))
soil.curve(to: NSPoint(x: 436, y: 291), controlPoint1: NSPoint(x: 494, y: 284), controlPoint2: NSPoint(x: 462, y: 284))
soil.curve(to: NSPoint(x: 380, y: 262), controlPoint1: NSPoint(x: 418, y: 299), controlPoint2: NSPoint(x: 392, y: 286))
soil.close()
let soilGradient = NSGradient(colors: [
    color(116, 98, 72),
    color(74, 61, 48),
    color(43, 36, 29),
])!
soilGradient.draw(in: soil, angle: 90)

let soilHighlight = NSBezierPath()
soilHighlight.move(to: NSPoint(x: 437, y: 293))
soilHighlight.curve(to: NSPoint(x: 582, y: 288), controlPoint1: NSPoint(x: 474, y: 304), controlPoint2: NSPoint(x: 544, y: 301))
soilHighlight.lineWidth = 6
soilHighlight.lineCapStyle = .round
color(188, 165, 123, 0.28).setStroke()
soilHighlight.stroke()

image.unlockFocus()

guard
    let tiffData = image.tiffRepresentation,
    let bitmap = NSBitmapImageRep(data: tiffData),
    let pngData = bitmap.representation(using: .png, properties: [:])
else {
    fputs("failed to generate icon data\n", stderr)
    exit(1)
}

do {
    try pngData.write(to: URL(fileURLWithPath: outputPath))
} catch {
    fputs("failed to write icon: \(error)\n", stderr)
    exit(1)
}
