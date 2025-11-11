#[derive(Debug)]
pub struct Tailer {
    reader: tokio::io::BufReader<tokio::fs::File>,
}

impl Tailer {
    pub async fn new(file_path: &str) -> std::io::Result<Self> {
        use tokio::io::AsyncSeekExt;

        let file = tokio::fs::File::open(file_path).await?;
        let mut reader = tokio::io::BufReader::new(file);

        reader.seek(SeekFrom::End(0)).await?;

        Ok(Tailer { reader })
    }

    // 현재 Seek 위치부터 추가된 행이 더 있다면 최대 10행까지 읽어서 반환하고, Seek 위치를 업데이트합니다.
    pub async fn tail(&mut self, num_lines: usize) -> Result<Vec<String>, std::io::Error> {
        use tokio::io::AsyncBufReadExt;

        let mut lines = vec![];

        // 현재 Seek 위치부터 읽기 시작
        let mut line = String::new();
        while lines.len() < num_lines && self.reader.read_line(&mut line).await? > 0 {
            lines.push(line.trim().to_string());
            line.clear();
        }

        Ok(lines)
    }
}
