package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"tritan.dev/bgp-tool/commands"
	"tritan.dev/bgp-tool/regex"
)

func main() {
	app := fiber.New()

	endpoints := []string{
		"/show-route?ip=<ip>",
		"/asn-routes?asn=<asn>",
		"/ping?ip=<ip>",
		"/traceroute?ip=<ip>",
	}

	app.Get("/", func(c *fiber.Ctx) error {
		endpointList := fmt.Sprintf("suck yourself ~as393577~\n\nAvailable Endpoints:\n\n%s", strings.Join(endpoints, "\n"))
		return c.SendString(endpointList)
	})

	app.Get("/show-route", func(c *fiber.Ctx) error {
		ip := c.Query("ip")
		if !regex.IsValidSubnet(ip) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid subnet param")
		}
		response, err := commands.ExecuteBirdCommand(ip)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("Route Info for IP %s:\n%s", ip, response))
	})

	app.Get("/asn-routes", func(c *fiber.Ctx) error {
		asn := c.Query("asn")
		if !regex.IsValidASN(asn) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ASN param")
		}
		response, err := commands.ExecuteBirdCommand(fmt.Sprintf("where bgp_path ~ [= * %s * =] all", asn))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("BGP Routes for ASN %s:\n%s", asn, response))
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		ip := c.Query("ip")
		if !regex.IsValidIP(ip) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid IP param")
		}
		response, err := commands.ExecutePing(ip)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("Ping for IP %s:\n%s", ip, response))
	})

	app.Get("/traceroute", func(c *fiber.Ctx) error {
		ip := c.Query("ip")
		if !regex.IsValidIP(ip) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid IP param")
		}
		response, err := commands.ExecuteTraceroute(ip)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("Traceroute for IP %s:\n%s", ip, response))
	})

	log.Fatal(app.Listen(":4000"))
}