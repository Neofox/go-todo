const path = require("path");
const { RspackManifestPlugin } = require("rspack-manifest-plugin");

/** @type {import('@rspack/cli').Configuration} */
const config = {
    entry: {
        main: "./static/js/index.ts",
    },
    output: {
        path: path.resolve(__dirname, "static/build"),
        clean: true,
        filename: "[name].[contenthash].js",
        chunkFilename: "chunks/[name].[contenthash].js",
    },
    resolve: {
        alias: {
            react: "preact/compat",
            "react-dom/test-utils": "preact/test-utils",
            "react-dom": "preact/compat",
            "react/jsx-runtime": "preact/jsx-runtime",
        },
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    {
                        loader: "postcss-loader",
                        options: {
                            postcssOptions: {
                                plugins: {
                                    tailwindcss: {},
                                    autoprefixer: {},
                                },
                            },
                        },
                    },
                ],
                type: "css",
            },
            {
                test: /\.tsx?$/,
                use: {
                    loader: "builtin:swc-loader",
                    options: {
                        jsc: {
                            parser: {
                                syntax: "typescript",
                                tsx: true,
                            },
                        },
                    },
                },
                type: "javascript/auto",
            },
        ],
    },
    resolve: {
        extensions: [".ts", ".tsx", ".js", ".jsx"],
    },
    experiments: {
        css: true,
    },
    devServer: {
        port: 8000,
        open: false,
        proxy: [
            {
                context: ["/static/build"],
                target: "http://localhost:7331",
            },
        ],
        devMiddleware: {
            writeToDisk: true,
            index: true,
            publicPath: "/static/build",
        },
    },
    optimization: {
        splitChunks: {
            chunks: "async",
            minSize: 1,
            cacheGroups: {
                vendor: {
                    test: /[\\/]node_modules[\\/]/,
                    name: "vendors",
                    chunks: "all",
                },
            },
        },
        runtimeChunk: "single",
    },
    plugins: [
        new RspackManifestPlugin({
            fileName: "manifest.json",
            publicPath: "/static/build",
        }),
    ],
};

module.exports = config;
