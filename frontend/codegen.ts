import type { CodegenConfig } from "@graphql-codegen/cli";

// Reads the gqlgen schema-first SDL and generates a typed `graphql()` function
// plus operation types into src/gql/. Used with urql via
// `useQuery({ query: graphql(`...`) })` for fully typed documents.
const config: CodegenConfig = {
  schema: "../graph/*.graphqls",
  documents: ["src/**/*.{ts,tsx}", "!src/gql/**/*"],
  ignoreNoDocuments: true,
  generates: {
    "src/gql/": {
      preset: "client",
      config: {
        useTypeImports: true,
      },
    },
  },
};

export default config;
