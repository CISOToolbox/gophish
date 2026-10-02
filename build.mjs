/*
 * build.mjs
 *
 * Front-end build for Gophish, powered by esbuild.
 * Replaces the previous gulp + webpack + babel toolchain:
 *   - vendor libraries are concatenated (order matters) and minified
 *   - plain, non-module app scripts are minified in place (globals preserved
 *     so inline HTML handlers keep working)
 *   - the ES module entries are bundled (format: iife)
 *   - stylesheets are concatenated and minified
 */

import * as esbuild from 'esbuild'
import { readFile, writeFile, mkdir } from 'node:fs/promises'
import path from 'node:path'

const vendorDir = 'static/js/src/vendor'
const appDir = 'static/js/src/app'
const cssDir = 'static/css'
const distJs = 'static/js/dist'
const distApp = 'static/js/dist/app'
const distCss = 'static/css/dist'

// Vendor libraries, concatenated in this exact order (jQuery first, plugins
// after, etc.). These are global-scope scripts, not ES modules.
const vendorFiles = [
  'jquery.js',
  'bootstrap.min.js',
  'moment.min.js',
  'papaparse.min.js',
  'd3.min.js',
  'topojson.min.js',
  'datamaps.min.js',
  'jquery.dataTables.min.js',
  'dataTables.bootstrap.js',
  'datetime-moment.js',
  'jquery.ui.widget.js',
  'jquery.fileupload.js',
  'jquery.iframe-transport.js',
  'sweetalert2.min.js',
  'bootstrap-datetime.js',
  'select2.min.js',
  'core.min.js',
  'highcharts.js',
  'ua-parser.min.js',
]

// Stylesheets, concatenated in this exact order.
const cssFiles = [
  'bootstrap.min.css',
  'main.css',
  'dashboard.css',
  'flat-ui.css',
  'dataTables.bootstrap.css',
  'font-awesome.min.css',
  'chartist.min.css',
  'bootstrap-datetime.css',
  'checkbox.css',
  'sweetalert2.min.css',
  'select2.min.css',
  'select2-bootstrap.min.css',
]

// ES module entries that must be bundled (resolve imports like zxcvbn).
const bundledEntries = ['passwords', 'users', 'webhooks']

// Plain, non-module app scripts: minified individually, kept in global scope.
const plainScripts = [
  'autocomplete',
  'campaign_results',
  'campaigns',
  'dashboard',
  'groups',
  'landing_pages',
  'sending_profiles',
  'settings',
  'templates',
  'gophish',
]

async function concatMinifyJs(files, dir, outFile) {
  const sources = await Promise.all(
    files.map((f) => readFile(path.join(dir, f), 'utf8')),
  )
  const { code } = await esbuild.transform(sources.join('\n;\n'), {
    minify: true,
    legalComments: 'inline',
  })
  await writeFile(outFile, code)
}

async function concatMinifyCss(files, dir, outFile) {
  const sources = await Promise.all(
    files.map((f) => readFile(path.join(dir, f), 'utf8')),
  )
  const { code } = await esbuild.transform(sources.join('\n'), {
    loader: 'css',
    minify: true,
    legalComments: 'inline',
  })
  await writeFile(outFile, code)
}

async function minifyPlainScripts(names, dir, outDir) {
  await Promise.all(
    names.map(async (name) => {
      const src = await readFile(path.join(dir, `${name}.js`), 'utf8')
      const { code } = await esbuild.transform(src, {
        minify: true,
        legalComments: 'inline',
      })
      await writeFile(path.join(outDir, `${name}.min.js`), code)
    }),
  )
}

async function bundleEntries(names, dir, outDir) {
  await esbuild.build({
    entryPoints: Object.fromEntries(
      names.map((n) => [n, path.join(dir, `${n}.js`)]),
    ),
    outdir: outDir,
    entryNames: '[name].min',
    bundle: true,
    minify: true,
    format: 'iife',
    target: 'es2015',
    legalComments: 'inline',
    logLevel: 'info',
  })
}

async function main() {
  await Promise.all([
    mkdir(distApp, { recursive: true }),
    mkdir(distCss, { recursive: true }),
  ])

  await Promise.all([
    concatMinifyJs(vendorFiles, vendorDir, path.join(distJs, 'vendor.min.js')),
    concatMinifyCss(cssFiles, cssDir, path.join(distCss, 'gophish.css')),
    minifyPlainScripts(plainScripts, appDir, distApp),
    bundleEntries(bundledEntries, appDir, distApp),
  ])

  console.log('Front-end build complete.')
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
