extern crate log4rs;
use log::{debug, error, info, trace, warn};

fn main() {
    log4rs::init_file("log4rs.yml", Default::default()).unwrap();

    error!("asdf");
}
