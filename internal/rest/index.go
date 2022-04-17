package rest

import (
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", htmlContentLength)
	w.Write(html)
}

// TODO: Consider Go 1.16 embed
var html = []byte(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Go PDF to Image</title>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
			}
			form {
				display: table;
			}
			p {
				display: table-row;
			}
			label {
				display: table-cell;
				margin: 4px;
				text-align: end;
			}
			input, select {
				display: table-cell;
				margin: 4px;
			}
			.main {
				width: 360px;
				margin: 0 auto;
			}
		</style>
		<script>
			function update() {
				var p = ["width", "height", "format", "quality", "resize", "resample", "pages", "archive"]
					.map(function(id) {
						var v = document.getElementById(id).value;
						return v && v !== "" ? id + "=" + v : null;
					}).filter(function(v) { return v; }).join("&");
				var url = "scale?" + p;
				document.getElementById("url").textContent = url;
				document.getElementById("form1").action = url;
			}
		</script>
		<noscript>
			Please enable JavaScript
		</noscript>
	</head>
	<body onload="update()">
		<div class="main">
			<h1>Convert PDF Document</h1>
			<form id="form1" action="scale" method="POST" enctype="multipart/form-data">
				<p>
					<label for="convert">Document*:</label>
					<input id="convert" type="file" name="pdf" accept="application/pdf">
				</p>
				<p>
					<label for="width">Width*:</label>
					<input id="width" type="number" min="1" max="65500" value="256" onchange="update()">
				</p>
				<p>
					<label for="height">Height*:</label>
					<input id="height" type="number" min="1" max="65500" value="256" onchange="update()">
				</p>
				<p>
					<label for="format">Format*:</label>
					<select id="format" onchange="update()">
						<option value="jpeg">JPEG</option>
						<option value="png">PNG</option>
						<option value="bmp">BMP</option>
						<option value="gif">GIF</option>
						<option value="tiff">TIFF</option>
					</select>
				</p>
				<p>
					<label for="quality">Quality<sup>1</sup>:</label>
					<input id="quality" type="number" min="0" max="100" placeholder="99" onchange="update()">
				</p>
				<p>
					<label for="resize">Resizing:</label>
					<select id="resize" onchange="update()">
						<option>
						<optgroup label="Fit">
							<option value="fit">Fit (w/o upscaling, no bars)</option>
							<option value="fit-black">Fit (w/o upscaling, black bars)</option>
							<option value="fit-white">Fit (w/o upscaling, white bars)</option>
							<option value="fit-upscale">Fit (w/ upscaling, no bars)</option>
							<option value="fit-upscale-black">Fit (w/ upscaling, black bars)</option>
							<option value="fit-upscale-white">Fit (w/ upscaling, white bars)</option>
						</optgroup>
						<optgroup label="Fill">
							<option value="fill-top-left">Fill (top left)</option>
							<option value="fill-top">Fill (top center)</option>
							<option value="fill-top-right">Fill (top right)</option>
							<option value="fill-left">Fill (middle left)</option>
							<option value="fill">Fill (middle center)</option>
							<option value="fill-right">Fill (middle right)</option>
							<option value="fill-bottom-left">Fill (bottom left)</option>
							<option value="fill-bottom">Fill (bottom center)</option>
							<option value="fill-bottom-right">Fill (bottom right)</option>
						</optgroup>
						<optgroup label="Other">
							<option value="stretch">Stretch</option>
						</optgroup>
					</select>
				</p>
				<p>
					<label for="resample">Resampling:</label>
					<select id="resample" onchange="update()">
						<option>
						<option value="box">Box sampling</option>
						<option value="nearest">Nearest neighbor</option>
						<option value="linear">Linear resampling</option>
						<option value="hermite">Hermite interpolation</option>
						<option value="mitchell">Mitchell-Netravali</option>
						<option value="catmull">Catmull-Rom</option>
						<option value="bspline">B-spline (B=1, C=0)</option>
						<option value="bartlett">Bartlett window sinc</option>
						<option value="lanczos">Lanczos</option>
						<option value="hann">Hann window sinc</option>
						<option value="hamming">Hamming window sinc</option>
						<option value="blackman">Blackman window sinc</option>
						<option value="welch">Welch parabolic window sinc</option>
						<option value="cosine">Cosine window sinc</option>
					</select>
				</p>
				<p>
					<label for="pages">Pages<sup>2</sup>:</label>
					<input id="pages" type="text" pattern="^(\d{0,7})(?:\-(_|\d{1,7}))?$" placeholder="1-_" onchange="update()">
				</p>
				<p>
					<label for="archive">Archive<sup>3</sup>:</label>
					<select id="archive" onchange="update()">
						<option>
						<option value="zip">ZIP</option>
						<option value="tar">Tarball</option>
					</select>
				</p>
				<p>
					<label for="upload">Upload:</label>
					<input id="upload" type="submit">
				</p>
			</form>
			<p>
				Resulting URL query:
				<textarea id="url" rows="3" cols="40" readonly></textarea>
			</p>
			<h1>Retrieve Page Count</h1>
			<form action="scale" method="POST" enctype="multipart/form-data">
				<p>
					<label for="count">Document*:</label>
					<input id="count" type="file" name="pdf" accept="application/pdf">
				</p>
				<p>
					<label for="upload">Upload:</label>
					<input id="upload" type="submit">
				</p>
			</form>
			<h1>Notes</h1>
			<p>* Mandatory field</p>
			<p><sup>1</sup> 0-100 for JPEGs, 0-9 for PNG compression</p>
			<p><sup>2</sup> For instance, 1 or 1-_ or 1-10</p>
			<p><sup>3</sup> Mandatory if pages is set to multiple</p>
		</div>
	</body>
</html>
`)

var htmlContentLength = strconv.Itoa(len(html))
