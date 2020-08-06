var reporter = require('cucumber-html-reporter');

if (process.argv[2]) {
        jsonDir = process.argv[2]
} else {
        throw TypeError('Probr HTML Generator expects exactly one argument.')
}

var options = {
        theme: 'bootstrap',                        
        jsonDir: jsonDir,
        output: `${jsonDir}/cucumber_report.html`,
        reportSuiteAsScenarios: true,
        scenarioTimestamp: true,
        launchReport: true,
        name: "Probr - Compliance as Code",
        brandTitle: "Citihub Digital",
        metadata: {
            "App Version":"0.3.2",
            "Test Environment": "DEV",            
            "Platform": "Windows 10",
            "Parallel": "Scenarios",
            "Executed": "Remote"
        }
    };

    reporter.generate(options);
    