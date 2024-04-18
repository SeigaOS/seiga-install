use std::{fs::read_to_string, path::PathBuf, str::FromStr};

use proc_macro::TokenStream;
use quote::quote;

#[proc_macro]
pub fn embed_schema(_: TokenStream) -> TokenStream {
    let mut schem = PathBuf::from_str(env!("CARGO_MANIFEST_DIR")).unwrap();
    schem.pop();
    schem.push("interface");
    schem.push("schema.toml");
    let schem = read_to_string(schem).expect("broken repository files");
    quote! {#schem}.into()
}