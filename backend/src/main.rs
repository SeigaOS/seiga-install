use std::io;
use std::fs::File;
use std::io::prelude::*;
use crossterm::event::{self, Event, KeyCode, KeyEvent, KeyEventKind};
use ratatui::{
    prelude::*,
    symbols::border,
    widgets::{block::*, *},
};

mod tui;

fn main() -> io::Result<()> {
    
    let mut terminal = tui::init()?;
    let app_result = App::default().run(&mut terminal);
    tui::restore()?;
    app_result
}

#[derive(Debug, Default)]
pub struct App {
    countertime: u16,
    countercel: u16,
    pos: bool,
    exit: bool,
}

//true = time false = cel

impl App {
    pub fn run(&mut self, terminal: &mut tui::Tui) -> io::Result<()> {
        self.countertime = 10;
        self.countercel = 30;

        while !self.exit {
            terminal.draw(|frame| self.render_frame(frame))?;
            self.handle_events()?;
        }
        Ok(())
    }

    fn render_frame(&self, frame: &mut Frame) {
        frame.render_widget(self, frame.size());
    }

    fn handle_events(&mut self) -> io::Result<()> {
        match event::read()? {
            Event::Key(key_event) if key_event.kind == KeyEventKind::Press => {
                self.handle_key_event(key_event)
            }
            _ => {}
        };
        Ok(())
    }

    fn handle_key_event(&mut self, key_event: KeyEvent) {
        match key_event.code {
            KeyCode::Char('q') => self.exit(),
            KeyCode::Left => self.decrement_counter(),
            KeyCode::Right => self.increment_counter(),
            KeyCode::Down => self.uppos(),
            KeyCode::Up => self.downpos(),
            KeyCode::Enter => self.makefolder().expect("WORKED<3"),
            _ => {}
        }
    }

    fn exit(&mut self) {
        self.exit = true;
    }
    fn makefolder(&mut self) -> std::io::Result<()> {
        let mut file = File::create("newsettings.toml")?;
        write!(file, "[tempcheck]\nchecktimems = {} \ncheckthresholdcelc = {} ", self.countertime, self.countercel)?;
        self.exit();
        Ok(())
    }
    //false == down true == up
    fn downpos(&mut self) {
        if self.pos == false {
            self.pos = true;
        } else {
            self.pos = false;
        }
    }
    fn uppos(&mut self) {
        if self.pos == true {
            self.pos = false;
        } else {
            self.pos = true;
        }
    }
    fn increment_counter(&mut self) {
        if self.pos == false {
            if self.countertime >= 10 {
                self.countertime += 10;
            }
        } else {
            if self.countercel >= 30 {
                self.countercel += 1;
            }
        }
    }

    fn decrement_counter(&mut self) {
        if self.pos == false {
            if self.countertime >= 100 {
                self.countertime -= 10;
            }
        } else {
            if self.countercel >= 30 {
                self.countercel += 1;
            }        
        }
    }
}

impl Widget for &App {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let mut topac = "";
        let mut botomac = "";

        let title = Title::from(" System's Setting Editor ".bold());
        let instructions = Title::from(Line::from(vec![
            " Decrement ".into(),
            "<Left>".blue().bold(),
            " Increment ".into(),
            "<Right>".blue().bold(),
            " down ".into(),
            "<Down>".blue().bold(),
            " up ".into(),
            "<Up>".blue().bold(),
            " enter ".into(),
            "<Enter>".blue().bold(),
            " Quit ".into(),
            "<Q> ".blue().bold(),
        ]));
        let block = Block::default()
            .title(title.alignment(Alignment::Center))
            .title(
                instructions
                    .alignment(Alignment::Center)
                    .position(Position::Bottom),
            )
            .title("secondig")
            .borders(Borders::ALL)
            .border_set(border::THICK);
        if self.pos == false { 
            topac = ">";
            botomac = "";
        } else {
            topac = "";
            botomac = ">";
        }
        let counter_time = Text::from(vec![Line::from(vec![
            format!("{} ", topac).into(),
            "CPU check time interval: ".into(),
            self.countertime.to_string().yellow(),
            " seconds".blue().into(),
       
        ])]);
        let counter_temp = Text::from(vec![Line::from(vec![
            format!("{} ", botomac).into(),
            "CPU tempurate threshold: ".into(),
            self.countercel.to_string().yellow(),
            " Celcius".blue().into(),
        ])]);

        let blocknew: Block = Block::new()
            .title("hey <3 system settings editor <3 ")
            .padding(Padding::new(5, 10, 1, 2));

        Paragraph::new(counter_time)
            .centered()
            .block(block)
            .render(area, buf);
    
        Paragraph::new(counter_temp)
                .centered()
                .block(blocknew)
                .render(area, buf);
        }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn render() {
        let app = App::default();
        let mut buf = Buffer::empty(Rect::new(0, 0, 50, 4));

        app.render(buf.area, &mut buf);

        let mut expected = Buffer::with_lines(vec![
            "┏━━━━━━━━━━━━━ System's settings editor<3 ━━━━━━━━━━━━━━━━━━━━━┓",
            "┃           > CPU check time interval: 10 seconds              ┃",
            "┃             CPU tempurate threshold: 30 Celcius              ┃",
            "┗━ Decrement <Left> Increment <Right> Enter <enter> Quit <Q> ━━┛",
        ]);
        let title_style = Style::new().bold();
        let counter_style = Style::new().yellow();
        let key_style = Style::new().blue().bold();
        expected.set_style(Rect::new(14, 0, 22, 1), title_style);
        expected.set_style(Rect::new(28, 1, 1, 1), counter_style);
        expected.set_style(Rect::new(13, 3, 6, 1), key_style);
        expected.set_style(Rect::new(30, 3, 7, 1), key_style);
        expected.set_style(Rect::new(43, 3, 4, 1), key_style);
        assert_eq!(buf, expected);
    }



}
