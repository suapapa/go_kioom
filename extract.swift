import Foundation
import Quartz

let arguments = CommandLine.arguments
if arguments.count < 2 {
    print("Usage: extract_pdf <pdf_path>")
    exit(1)
}

let pdfPath = arguments[1]
let url = URL(fileURLWithPath: pdfPath)
guard let document = PDFDocument(url: url) else {
    print("Could not load PDF")
    exit(1)
}

var fullText = ""
for i in 0..<document.pageCount {
    if let page = document.page(at: i) {
        if let string = page.string {
            fullText += "--- Page \(i+1) ---\n"
            fullText += string
            fullText += "\n\n"
        }
    }
}

let outputPath = arguments.count > 2 ? arguments[2] : nil

if let outBase = outputPath {
    do {
        try fullText.write(toFile: outBase, atomically: true, encoding: .utf8)
        print("Successfully written to \(outBase)")
    } catch {
        print("Failed to write file: \(error)")
    }
} else {
    print(fullText)
}
