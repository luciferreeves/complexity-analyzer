# Algorithm Complexity Analyzer

A real-time algorithm complexity analyzer that uses AI to determine Big O notation and generates performance graphs through dynamic execution testing.

![Algorithm Complexity Analyzer](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![Gemini AI](https://img.shields.io/badge/Gemini-2.5%20Flash-4285F4?style=flat&logo=google)
![License](https://img.shields.io/badge/license-MIT-green)

## 🚀 Features

- **AI-Powered Analysis**: Uses Google's Gemini 2.5 Flash to analyze algorithm complexity
- **Multi-Language Support**: Analyze algorithms written in JavaScript, Python, Go, and more
- **Real-Time Performance Testing**: Dynamically executes code with varying input sizes
- **Adaptive Benchmarking**: Automatically adjusts iterations to maintain optimal test duration
- **Interactive Visualization**: Real-time performance graphs using Chart.js
- **Static Code Analysis**: Detailed breakdown of algorithm structure and patterns

## 🔧 Installation

1. **Clone the repository:**

```bash
git clone https://github.com/luciferreeves/complexity-analyzer.git
cd complexity-analyzer
```

2. **Install dependencies:**

```bash
go mod download
```

3. **Set up environment variables:**

```bash
cp .env.example .env
```

Edit `.env` and add your Gemini API key:

```env
GEMINI_API_KEY=your-api-key-here
PORT=3000
```

4. **Run the application:**

```bash
go run main.go
```

5. **Open your browser:**

```
http://localhost:3000
```

## 🧠 How It Works

1. **AI Analysis**: Gemini AI analyzes your code and determines:

   - Time complexity notation
   - Algorithm approach and patterns
   - Appropriate test sizes

2. **Dynamic Testing**: The system:

   - Generates complete Go test program
   - Runs algorithm with increasing input sizes
   - Adaptively adjusts iterations
   - Measures execution time with high precision

3. **Visualization**: Results are displayed with:
   - Complexity badge with confidence score
   - Static analysis insights
   - Interactive performance graph
   - Detailed timing breakdown

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ using Go**
