const path = require("path");

/** @type {import('@rspack/cli').Configuration} */
const config = {
    entry: {
        main: "./static/js/index.ts",
    },
    output: {
        path: path.resolve(__dirname, "static/build"),
        clean: true,
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
};

module.exports = config;
