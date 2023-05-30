package engine

import (
	"testing"

	"github.com/ragrag/envoy/pkg/infra"
	"github.com/ragrag/envoy/pkg/mock"
	"github.com/ragrag/envoy/pkg/runtimeenv"
	"github.com/ragrag/envoy/pkg/sandbox"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

var c = CreateTestContainer()

func CreateTestContainer() *dig.Container {
	c := dig.New()

	c.Provide(func() *infra.Config {
		return mock.ConfigMock
	})
	c.Provide(infra.NewLogger)
	c.Provide(runtimeenv.NewRuntimeProvider)
	c.Provide(sandbox.NewSanboxManager)
	c.Provide(NewEngine)

	err := c.Invoke(func(logger *logrus.Logger, config *infra.Config, runtimeProvider *runtimeenv.RuntimeProvider, engine *Engine) error {
		e := runtimeProvider.Load()
		if e != nil {
			return e
		}

		return engine.Ignite()
	})

	if err != nil {
		panic("Failed to create test container")
	}

	return c
}

func testHelloWorld(t *testing.T, language string, code string) {
	err := c.Invoke(func(engine *Engine) error {
		params := &RunParams{
			Language: language,
			Code:     code,
		}

		res, err := engine.RunCode(params)

		if err != nil {
			return err
		}

		if res.Status != SUCCESS {
			t.Errorf("\nExpected Status:\n%s\nGot:\n%s", SUCCESS, res.Status)
		}
		if res.Output != "El Psy Kongroo\n" {
			t.Errorf("\nExpected Output:\n%s\nGot:\n%s", "El Psy Kongroo\n", res.Output)
		}

		return nil
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGo(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "go", `
		package main
		import "fmt"
		func main() {
			fmt.Println("El Psy Kongroo")
		}
	`)
}

func TestRust(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "rust", `
		fn main() {
			println!("El Psy Kongroo");
		}
	`)
}

func TestJavaScript(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "javascript", `
	console.log("El Psy Kongroo");
	`)
}

func TestTypeScript(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "typescript", `
	console.log("El Psy Kongroo");
	`)
}

func TestCPP(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "cpp", `
	#include <iostream>
	using namespace std;

	int main() {
		cout << "El Psy Kongroo" << endl;
		return 0;
	}
	`)
}

func TestC(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "c", `
	#include <stdio.h>

	int main() {
		printf("El Psy Kongroo\n");
		return 0;
	}
	`)
}

func TestJava(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "java", `
		public class ElPsyKongroo {
			public static void main(String[] args) {
				System.out.println("El Psy Kongroo");
			}
		}
	`)
}

func TestPython2(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "python2", `print("El Psy Kongroo")`)
}

func TestPython3(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "python3", `print("El Psy Kongroo")`)
}

func TestCSharp(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "csharp", `
	using System;

	public class Test
	{
		public static void Main(string[] args)
		{
			Console.WriteLine("El Psy Kongroo");
		}
	}
	`)
}

func TestPHP(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "php", `<?php echo "El Psy Kongroo\n"; ?>`)
}

func TestSwift(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "swift", `print("El Psy Kongroo")`)
}

func TestKotlin(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "kotlin", `
	fun main() {
		println("El Psy Kongroo")
	}
	`)
}

func TestRuby(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "ruby", `puts("El Psy Kongroo")`)
}

func TestScala(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "scala", `
		object Main {
			def main(args: Array[String]) = {
				println("El Psy Kongroo")
			}
		}
	`)
}

func TestErlang(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "erlang", `
	main(_) ->
		io:fwrite("El Psy Kongroo\n").
	`)
}

func TestElixir(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "elixir", `IO.puts "El Psy Kongroo"`)
}

func TestHaskell(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "haskell", `
		main = putStrLn "El Psy Kongroo"
	`)
}

func TestZig(t *testing.T) {
	t.Parallel()
	testHelloWorld(t, "zig", `
		const std = @import("std");
		
		pub fn main() !void {
			const stdout = std.io.getStdOut().writer();
			try stdout.print("El Psy Kongroo\n", .{});
		}
	`)
}
