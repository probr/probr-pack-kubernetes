var reporter = require('cucumber-html-reporter');

var options = {
        theme: 'bootstrap',                        
        jsonDir: `${process.argv[2]}`,
        output: `${process.argv[2]}/cucumber_report.html`,
        reportSuiteAsScenarios: true,
        scenarioTimestamp: true,
        launchReport: true,
        metadata: {
            "App Version":"0.3.2",
            "Test Environment": "DEV",            
            "Platform": "Windows 10",
            "Parallel": "Scenarios",
            "Executed": "Remote"
        }
    };

    reporter.generate(options);
    