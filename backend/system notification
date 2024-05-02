

use sysinfo::{System, Components};
use std::{process::Command, string};


fn cputils() {
    let components = Components::new_with_refreshed_list();
    println!("=> components:");
    for component in &components {
        let component_name = component.label();
        if component_name.contains("tdie") {
            // println!("{} with temp {:?}C", component_name, component.temperature());
            // println!("{}", component_name);
            if component_name.contains("tdie6") {
                println!("{}", component.temperature());
                if component.temperature() > 39.0 {
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
    let mut sys = System::new_all();
    sys.refresh_all();
    cputils();
    battman();
    // println!("total memory: {} bytes", sys.total_memory());
    // println!("used memory : {} bytes", sys.used_memory());
    // println!("total swap  : {} bytes", sys.total_swap());
    // println!("used swap   : {} bytes", sys.used_swap());
    // println!("System name:             {:?}", System::name());
    // println!("System kernel version:   {:?}", System::kernel_version());
    // println!("System OS version:       {:?}", System::os_version());
    // println!("System host name:        {:?}", System::host_name());
    
    // loop {
    //     sys.refresh_cpu();
    //     for cpu in sys.cpus() {
    //         print!("{}% ", cpu.cpu_usage());
    //     }
    //     std::thread::sleep(sysinfo::MINIMUM_CPU_UPDATE_INTERVAL);
    // }
    
}
