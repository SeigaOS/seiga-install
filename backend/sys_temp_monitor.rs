

use sysinfo::{System, Components};
use std::{process::Command, string};
use std::{thread, time};
use serde_derive::Deserialize;
use std::fs;
use std::process::exit;
use toml;

#[derive(Deserialize)]
struct Data {
    tempcheck: Tempcheck,
}
#[derive(Deserialize)]
struct Tempcheck {
    checktimems: u64,
    checkthresholdcelc: f32,
}

fn cputils(temp: f32) {

    let components = Components::new_with_refreshed_list();
    println!("=> components:");
    for component in &components {
        let component_name = component.label();
        if component_name.contains("tdie") {
            println!("{} with temp {:?}C", component_name, component.temperature());
            // println!("{}", component_name);
            if component_name.contains("tdie3") {
                println!("{}", component.temperature());
                if component.temperature() > temp {
                    let stringfr = format!("display notification \"Your CPU temperature is too high\" with title \"CPU Temp Checker X 900\"");
                    Command::new("osascript")
                        .args(["-e", &stringfr])
                        .output() 
                        .expect("Failed to execute process");
                
                    println!("{}", stringfr);    
                    println!("Too high");    
                    }

                
            }

        }
        //osascript -e 'display notification "your cpu tempurature is too high" with title "cpu temp checker x 900 | your cpu tempurature is too high"'
    }
}


fn battman() {
    let output = { Command::new("pmset")
        .arg("-g")
        .arg("batt")
        .output()
        .expect("failed to execute process")
    };
    let stringer = String::from_utf8_lossy(&output.stdout);
    let mut parts = stringer.split(")");
    let mut parts = parts.nth(1).expect("worked").split(" ");
    let battpercent = parts.nth(0).expect("WORKED").trim();
    let timeleftpercent = parts.nth(1).expect("WORKED").trim();
    println!("current percentage {}", battpercent);
    println!("time remaining {}", timeleftpercent);
// '    println!("=> disks:");
//     let disks = Disks::new_with_refreshed_list();
//     for disk in &disks {
//         println!("{disk:?}");
//     }'
}
fn main() {
   // println!("line");
    //pmset -g batt 
    //chmod +x setlimitcheck

    let filename = "put the path to a checkfor toml here";
    let contents = match fs::read_to_string(filename) {
        Ok(c) => c,
        Err(_) => {
            eprintln!("Could not read file `{}`", filename);
            exit(1);
        }
    };

    let data: Data = match toml::from_str(&contents) {
        Ok(d) => d,
        Err(_) => {
            eprintln!("Unable to load data from `{}`", filename);
            exit(1);
        }
    };

    let mut sys = System::new_all();
    sys.refresh_all();
    loop {
        let ten_millis = time::Duration::from_millis(data.tempcheck.checktimems);
        let now = time::Instant::now();
        thread::sleep(ten_millis);
        cputils(data.tempcheck.checkthresholdcelc );
    };
    // battman();
    // println!("total memory: {} bytes", sys.total_memory());
    // println!("used memory : {} bytes", sys.used_memory());
    // println!("total swap  : {} bytes", sys.total_swap());
    // println!("used swap   : {} bytes", sys.used_swap());
    // println!("System name:             {:?}", System::name());
    // println!("System kernel version:   {:?}", System::kernel_version());
    // println!("System OS version:       {:?}", System::os_version());
    // println!("System host name:        {:?}", System::host_name());
    
    // loop {
    //     sys.refresh_cpu(); // Refreshing CPU information.
    //     for cpu in sys.cpus() {
    //         print!("{}% ", cpu.cpu_usage());
    //     }
    //     std::thread::sleep(sysinfo::MINIMUM_CPU_UPDATE_INTERVAL);
    // }
    
}
