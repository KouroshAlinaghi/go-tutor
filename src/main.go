package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type proficiencyLevel = string

const (
    junior proficiencyLevel = "junior"
    expert proficiencyLevel = "expert"
    senior proficiencyLevel = "senior"
    teamLead proficiencyLevel = "team_lead"
)

const (
    employees_file_path = "./data/employees.csv"
    working_hours_file_path = "./data/working_hours.csv"
    teams_file_path = "./data/teams.csv"
    salary_configs_file_path = "./data/salary_configs.csv"
)

type employee struct {
    id int
    name string
    age int
    proficiencyLevel proficiencyLevel
    team *team
    workingHours [30][24]bool
    rawSalary int
    bonusAmount int
    taxAmount int
}

type team struct {
    id int
    headMember *employee
    bonusMinWorkingHours int
    members []*employee
}

type salaryConfig struct {
    baseSalary int
    salaryPerHour int
    salaryPerExtraHour int
    officialWorkingHours int
    taxPercentage int
}

type database struct {
    employees []*employee
    teams []*team
    salaryConfigs map[proficiencyLevel]*salaryConfig
}

func main() {
    db := NewDatabase()
    db.CalculateSalaries()
    
    ProcessInput(db)
}

func Shout(err error) {
    if err != nil {
        fmt.Println(err.Error())
    }
}

func NewDatabase() database {
    db := database{}
    db.salaryConfigs = make(map[string]*salaryConfig)

    err := db.readEmployeesFile()
    Shout(err)
    err = db.readWorkingHoursFile()
    Shout(err)
    err = db.readTeamsFile()
    Shout(err)
    err = db.readSalaryConfigsFile()
    Shout(err)

    return db
}

func ProcessInput(db database) {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        if !scanner.Scan() {
            continue
        }
        words := strings.Split(scanner.Text(), " ")
        switch words[0] {
        case "report_salaries":
            db.ReportSalaries()
        case "report_employee_salary":
            employeeId, _ := strconv.Atoi(words[1])
            db.ReportEmployeeSalary(employeeId)
        case "report_team_salary":
            teamId, _ := strconv.Atoi(words[1])
            db.ReportTeamSalary(teamId)
        case "report_total_hours_per_day":
            start_day, _ := strconv.Atoi(words[1])
            end_day, _ := strconv.Atoi(words[2])
            db.ReportTotalHoursPerDay(start_day, end_day)
        case "report_employee_per_hour":
            start_hour , _ := strconv.Atoi(words[1])
            end_hour, _ := strconv.Atoi(words[2])
            db.ReportEmployeePerHour(start_hour, end_hour)
        case "show_salary_config":
            db.ShowSalaryConfig(words[1])
        case "update_salary_parameters":
            newBase, _ := strconv.Atoi(words[2])
            newPerHour, _ := strconv.Atoi(words[3])
            newPerExtraHour, _ := strconv.Atoi(words[4])
            newOfficialHours, _ := strconv.Atoi(words[5])
            newTaxPerc, _ := strconv.Atoi(words[6])
            db.UpdateSalaryConfig(words[1], newBase, newPerHour, newPerExtraHour, newOfficialHours, newTaxPerc)
        }
    }
}

func (db database) ReportSalaries() {
    for _, emp := range db.employees {
        fmt.Printf("ID: %v\n", emp.id)
        fmt.Printf("Name: %v\n", emp.name)
        fmt.Printf("Total Working Hours: %v\n", emp.GetTotalWorkingHours())
        fmt.Printf("Total Earning: %v\n", emp.GetTotalEarning())
        fmt.Printf("---\n")
    }
}

func (db database) ReportEmployeeSalary(employeeId int) {
    employee, err := db.getEmployee(employeeId)
    if err != nil {
        fmt.Println("EMPLOYEE_NOT_FOUND")
        return
    }

    fmt.Printf("ID: %v\n", employee.id)
    fmt.Printf("Name: %v\n", employee.name)
    fmt.Printf("Age: %v\n", employee.age)
    fmt.Printf("Level: %v\n", employee.proficiencyLevel)
    if employee.team == nil {
        fmt.Printf("Team ID: N/A\n")
    } else {
        fmt.Printf("Team ID: %v\n", employee.team.id)
    }
    fmt.Printf("Total Working Hours: %v\n", employee.GetTotalWorkingHours())
    fmt.Printf("Absent Days: %v\n", employee.GetAbsentDays())
    fmt.Printf("Salary: %v\n", employee.rawSalary)
    fmt.Printf("Bonus: %v\n", employee.bonusAmount)
    fmt.Printf("Tax: %v\n", employee.taxAmount)
    fmt.Printf("Total Earning: %v\n", employee.GetTotalEarning())

}

func (db database) ReportTeamSalary(teamId int) {
    team, err := db.getTeam(teamId)
    if err != nil {
        fmt.Println("TEAM_NOT_FOUND")
        return
    }

    fmt.Printf("ID: %v\n", team.id)
    fmt.Printf("Head ID: %v\n", team.headMember.id)
    fmt.Printf("Head Name: %v\n", team.headMember.name)
    fmt.Printf("Team Total Working Hours: %v\n", team.GetTotalWorkingHours())
    fmt.Printf("Average Member Working Hours: %v\n", team.GetTotalWorkingHours()/len(team.members))

    fmt.Println("---")

    for _, member := range team.members {
        fmt.Printf("Member ID: %v\n", member.id)
        fmt.Printf("Total Earning: %v\n", member.GetTotalEarning())
    }

    fmt.Println("---")
}

func (db database) ReportTotalHoursPerDay(start_day, end_day int) {
    daysWorkingHours := []float32{}
    for day := start_day; day <= end_day; day++ {
        daysWorkingHours = append(daysWorkingHours, float32(db.getTotalWorkingHoursOnDay(day)))
        fmt.Printf("Day #%v: %v\n", day, db.getTotalWorkingHoursOnDay(day))
    }
    
    max := getMaxOfSlice(daysWorkingHours)
    min := getMinOfSlice(daysWorkingHours)

    fmt.Print("Day(s) with Max Working Hours: ")
    for day, value := range daysWorkingHours {
        if value == max {
            fmt.Printf("%v ", day+start_day)
        }
    }
    fmt.Print("\n")

    fmt.Print("Day(s) with Min Working Hours: ")
    for day, value := range daysWorkingHours {
        if value == min {
            fmt.Printf("%v ", day+start_day)
        }
    }
    fmt.Print("\n")

    fmt.Println("---")
}

func (db database) ReportEmployeePerHour(start_hour, end_hour int) {
    if start_hour >= end_hour || start_hour < 0 || end_hour > 24 {
        fmt.Println("INVALID_ARGUMENTS")
        return
    }
    workingEmployeeOnHour := []float32{}
    for hour := start_hour; hour < end_hour; hour++ {
        fmt.Printf("%v-%v: %v\n", hour, hour+1, db.getAvgEmployeesWorkingOnHour(hour))
        workingEmployeeOnHour = append(workingEmployeeOnHour, float32(db.getWorkingEmployeesOnHour(hour)))
    }

    fmt.Println("---")

    max := getMaxOfSlice(workingEmployeeOnHour)
    min := getMinOfSlice(workingEmployeeOnHour)

    fmt.Print("Period(s) with Max Working Employees: ")
    for hour, value := range workingEmployeeOnHour {
        if value == max {
            fmt.Printf("%v-%v ", hour+start_hour, hour+start_hour+1)
        }
    }
    fmt.Print("\n")

    fmt.Print("Period(s) with Min Working Employees: ")
    for hour, value := range workingEmployeeOnHour {
        if value == min {
            fmt.Printf("%v-%v ", hour+start_hour, hour+start_hour+1)
        }
    }
    fmt.Print("\n")
    
}

func (db database) ShowSalaryConfig(level proficiencyLevel) {
    config, found := db.salaryConfigs[level]
    if found {
        fmt.Printf("Base Salary: %v\n", config.baseSalary)
        fmt.Printf("Salary Per Hour: %v\n", config.salaryPerHour)
        fmt.Printf("Salary Per Extra Hour: %v\n", config.salaryPerExtraHour)
        fmt.Printf("Official Working Hours: %v\n", config.officialWorkingHours)
        fmt.Printf("Tax Percentage: %v\n", config.taxPercentage)
    } else {
        fmt.Println("INVALID_LEVEL")
        return
    }
}

func (db *database) UpdateSalaryConfig(level proficiencyLevel, newBase, newPerHour, newPerExtraHour, newOfficialHours, newTaxPerc int) {
    config, found := db.salaryConfigs[level]
    if found {
        config.baseSalary = newBase
        config.salaryPerHour = newPerHour
        config.salaryPerExtraHour = newPerExtraHour
        config.officialWorkingHours = newOfficialHours
        config.taxPercentage = newTaxPerc
    } else {
        fmt.Println("INVALID_LEVEL")
        return
    }

    fmt.Println("OK")
}

func (db *database) readEmployeesFile() error {
    lines, err := readFile(employees_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        id, _ := strconv.Atoi(rows[0])
        age, _ := strconv.Atoi(rows[2])
        
        db.addEmployee(id, rows[1], age, rows[3])
    }

    return err
}

func (db *database) readWorkingHoursFile() error {
    lines, err := readFile(working_hours_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        employeeId, _ := strconv.Atoi(rows[0])
        day, _ := strconv.Atoi(rows[1])
        workingInterval := strings.Split(rows[2], "-")
        startTime, _ := strconv.Atoi(workingInterval[0])
        endTime, _ := strconv.Atoi(workingInterval[1])
        
        db.setWorkingHours(employeeId, day, startTime, endTime)
    }

    return err
}

func (db *database) readTeamsFile() error {
    lines, err := readFile(teams_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        teamId, _ := strconv.Atoi(rows[0])
        headMemberId, _ := strconv.Atoi(rows[1])
        membersIds := StringsToInts(strings.Split(rows[2], "$"))
        bonusMinWorkingHours, _ := strconv.Atoi(rows[3])
        
        db.addTeam(teamId, headMemberId, membersIds, bonusMinWorkingHours)
    }

    return err
}

func (db *database) readSalaryConfigsFile() error {
    lines, err := readFile(salary_configs_file_path)
    if err != nil {
        return err
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        rows := strings.Split(line, ",")

        baseSalary, _ := strconv.Atoi(rows[1])
        salaryPerHour, _ := strconv.Atoi(rows[2])
        salaryPerExtraHour, _ := strconv.Atoi(rows[3])
        officialWorkingHours, _ := strconv.Atoi(rows[4])
        taxPercentage, _ := strconv.Atoi(rows[5])
        
        db.addSalaryConfig(rows[0], baseSalary, salaryPerHour, salaryPerExtraHour, officialWorkingHours, taxPercentage)
    }

    return err
}

func readFile(filepath string) ([]string, error) {
    var lines []string
    file, err := os.Open(filepath)
    if err != nil {
        return lines, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    if err := scanner.Err(); err != nil {
        return lines, err
    }

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, err
}

func (db *database) addEmployee(id int, name string, age int, level proficiencyLevel) {
    newEmployee := employee{id: id, name: name, age: age, proficiencyLevel: level}
    db.employees = append(db.employees, &newEmployee)
}

func (db *database) setWorkingHours(employeeId, day, startTime, endTime int) {
    employee, _ := db.getEmployee(employeeId)
    for hour := startTime; hour < endTime; hour++ {
        employee.workingHours[day-1][hour] = true
    }
}

func (db *database) addTeam(teamId, headMemberId int, membersIds []int, bonusMinWorkingHours int) {
    newTeam := team{id: teamId, bonusMinWorkingHours: bonusMinWorkingHours}

    headMember, _ := db.getEmployee(headMemberId)
    newTeam.headMember = headMember

    for _, memberId := range membersIds {
        emp, _ := db.getEmployee(memberId)
        emp.team = &newTeam
        newTeam.members = append(newTeam.members, emp)
    }

    db.teams = append(db.teams, &newTeam)
}

func (db *database) addSalaryConfig(level proficiencyLevel, baseSalary, salaryPerHour, salaryPerExtraHour, officialWorkingHours, taxPercentage int) {
    newSalaryConfig := salaryConfig{baseSalary: baseSalary, salaryPerHour: salaryPerHour, salaryPerExtraHour: salaryPerExtraHour, officialWorkingHours: officialWorkingHours, taxPercentage: taxPercentage}
    db.salaryConfigs[level] = &newSalaryConfig
}

func (db database) getEmployee(id int) (*employee, error) {
    for _, employee := range db.employees {
        if employee.id == id {
            return employee, nil
        }
    }

    return nil, errors.New("Employee Not Found")
}

func (db database) getTeam(id int) (*team, error) {
    for _, team := range db.teams {
        if team.id == id {
            return team, nil
        }
    }

    return nil, errors.New("Team Not Found")
}

func (db *database) CalculateSalaries() {
    for _, emp := range db.employees {
        salaryConfig := db.salaryConfigs[emp.proficiencyLevel]
        emp.rawSalary = salaryConfig.baseSalary

        if emp.GetTotalWorkingHours() > salaryConfig.officialWorkingHours {
            emp.rawSalary += salaryConfig.officialWorkingHours * salaryConfig.salaryPerHour
            emp.bonusAmount = (emp.GetTotalWorkingHours() - salaryConfig.officialWorkingHours) * salaryConfig.salaryPerExtraHour
        } else {
            emp.rawSalary += salaryConfig.salaryPerHour * emp.GetTotalWorkingHours()
        }

        emp.taxAmount = (emp.rawSalary + emp.bonusAmount) * salaryConfig.taxPercentage / 100
    }
}

func (db database) getTotalWorkingHoursOnDay(day int) int {
    res := 0
    for _, employee := range db.employees {
        for _, hour := range employee.workingHours[day] {
            if hour {
                res++
            }
        }
    }
    return res
}

func (db database) getWorkingEmployeesOnHour(hour int) int {
    cnt := 0
    for _, emp := range db.employees {
        for _, day := range emp.workingHours {
            if (day[hour]) {
                cnt++
            }
        }
    }
    return cnt
}

func (db database) getAvgEmployeesWorkingOnHour(starting_hour int) float32 {
    var cnt float32 = 0
    for day := 0; day < 30; day++ {
        for _, emp := range db.employees {
            if emp.workingHours[day][starting_hour-1] {
                cnt++
            }
        }
    }
    return cnt/30.0
}

func (emp employee) GetTotalWorkingHours() int {
    counter := 0
    for _, day := range emp.workingHours {
        for _, hour := range day {
            if hour {
                counter++
            }
        }
    }

    return counter
}

func (emp employee) GetTotalEarning() int {
    return emp.rawSalary + emp.bonusAmount - emp.taxAmount
}

func (emp employee) GetAbsentDays() int {
    absentDays := 0
    didWorkToday := false
    for _, day := range emp.workingHours {
        didWorkToday = false
        for _, hour := range day {
            if hour {
                didWorkToday = true
                break
            }
        }
        if !didWorkToday {
            absentDays++
        }
    }
    return absentDays
}

func (t team) GetTotalWorkingHours() int {
    res := 0
    for _, member := range t.members {
        res += member.GetTotalWorkingHours()
    }
    return res
}

func StringsToInts(strings []string) []int {
    res := []int{}
    for _, str := range strings {
        num, _ := strconv.Atoi(str)
        res = append(res, num)
    }
    return res
}

func getMaxOfSlice(slice []float32) float32 {
    if len(slice) == 0 {
        return 0
    }

    max := slice[0]
    for _, value := range slice {
        if value > max {
            max = value
        }
    }
    return max
}

func getMinOfSlice(slice []float32) float32 {
    if len(slice) == 0 {
        return 0
    }

    min := slice[0]
    for _, value := range slice {
        if value < min {
            min = value
        }
    }
    return min
}
